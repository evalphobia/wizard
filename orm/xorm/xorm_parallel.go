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
	slaves := xpr.orm.Slaves(cond.Table)
	var sessions []Session
	for _, slave := range slaves {
		s := slave.NewSession()
		for _, w := range cond.Where {
			s.And(w.Statement, w.Args...)
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

	// execute query
	var errList []error
	results := make(chan reflect.Value, len(slaves))
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
	for range slaves {
		v := <-results
		e.Set(reflect.AppendSlice(e, v.Elem()))
	}
	if len(errList) > 0 {
		return errors.NewErrParallelQuery(errList)
	}

	return nil
}
