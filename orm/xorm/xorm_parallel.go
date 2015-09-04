package xorm

import (
	"reflect"

	"github.com/evalphobia/wizard/errors"
)

// XormParallel supports concurrent query
type XormParallel struct {
	orm *Xorm
}

// FindParallel executes SELECT query to all of the shards
func (xpr *XormParallel) FindParallel(listPtr interface{}, table interface{}, where string, args ...interface{}) error {
	cond := NewFindCondition(table)
	cond.And(where, args...)
	return xpr.FindParallelByCondition(listPtr, cond)
}

// CreateSessionsWithCondition creates new sessions with conditional clause
func (xpr *XormParallel) CreateSessionsWithCondition(cond FindCondition) []Session {
	var sessions []Session
	slaves := xpr.orm.Slaves(cond.Table)
	for _, slave := range slaves {
		s := slave.NewSession()
		for _, w := range cond.Where {
			s.And(w.Statement, w.Args...)
		}
		for _, in := range cond.WhereIn {
			s.In(in.Statement, in.Args...)
		}
		for _, o := range cond.OrderBy {
			if o.OrderByDesc {
				s.Desc(o.Name)
			} else {
				s.Asc(o.Name)
			}
		}
		if cond.Limit > 0 {
			s.Limit(cond.Limit, cond.Offset)
		}
		sessions = append(sessions, s)
	}
	return sessions
}

// FindParallelByCondition executes SELECT query to all of the shards with conditions
func (xpr *XormParallel) FindParallelByCondition(listPtr interface{}, cond FindCondition) error {
	vt := reflect.TypeOf(listPtr)
	if vt.Kind() != reflect.Ptr {
		return errors.NewErrArgType("listPtr must be a pointer")
	}
	elem := vt.Elem()
	if elem.Kind() != reflect.Slice && elem.Kind() != reflect.Map {
		return errors.NewErrArgType("listPtr must be a pointer of slice or map")
	}

	// create session with the condition
	sessions := xpr.CreateSessionsWithCondition(cond)
	length := len(sessions)

	// execute query
	var errList []error
	results := make(chan reflect.Value, length)
	for _, s := range sessions {
		list := reflect.New(elem)
		go func(s Session, list reflect.Value) {
			err := s.Find(list.Interface())
			if err != nil {
				errList = append(errList, err)
			}
			results <- list
		}(s, list)
	}

	// wait for the results
	e := reflect.ValueOf(listPtr).Elem()
	for i := 0; i < length; i++ {
		v := <-results
		e.Set(reflect.AppendSlice(e, v.Elem()))
	}
	if len(errList) > 0 {
		return errors.NewErrParallelQuery(errList)
	}

	return nil
}

// CountParallelByCondition executes SELECT COUNT(*) query to all of the shards with conditions
func (xpr *XormParallel) CountParallelByCondition(objPtr interface{}, cond FindCondition) ([]int64, error) {
	vt := reflect.TypeOf(objPtr)
	if vt.Kind() != reflect.Ptr {
		return nil, errors.NewErrArgType("objPtr must be a pointer")
	}

	// create session with the condition
	sessions := xpr.CreateSessionsWithCondition(cond)
	length := len(sessions)

	// execute query
	var errList []error
	results := make(chan int64, length)
	for _, s := range sessions {
		var count int64
		go func(s Session, count int64) {
			count, err := s.Count(objPtr)
			if err != nil {
				errList = append(errList, err)
			}
			results <- count
		}(s, count)
	}

	if len(errList) > 0 {
		return nil, errors.NewErrParallelQuery(errList)
	}

	// wait for the results
	var counts []int64
	for i := 0; i < length; i++ {
		v := <-results
		counts = append(counts, v)
	}

	return counts, nil
}
