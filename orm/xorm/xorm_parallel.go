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
	sessions := xpr.CreateFindSessions(cond)
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
	sessions := xpr.CreateFindSessions(cond)
	length := len(sessions)

	// execute query
	var errList []error
	results := make(chan int64, length)
	for _, s := range sessions {
		go func(s Session) {
			count, err := s.Count(objPtr)
			if err != nil {
				errList = append(errList, err)
			}
			results <- count
		}(s)
	}

	// wait for the results
	var counts []int64
	for i := 0; i < length; i++ {
		v := <-results
		counts = append(counts, v)
	}
	if len(errList) > 0 {
		return counts, errors.NewErrParallelQuery(errList)
	}

	return counts, nil
}

// CreateFindSessions creates new sessions with conditional clause
func (xpr *XormParallel) CreateFindSessions(cond FindCondition) []Session {
	var sessions []Session
	slaves := xpr.orm.Slaves(cond.Table)

	for _, slave := range slaves {
		s := slave.NewSession()
		if len(cond.Columns) != 0 {
			s.Cols(cond.Columns...)
		}
		if cond.Selects != "" {
			s.Select(cond.Selects)
		}

		for _, w := range cond.Where {
			s.And(w.Statement, w.Args...)
		}
		for _, in := range cond.WhereIn {
			s.In(in.Statement, in.Args...)
		}
		for _, group := range cond.Group {
			s.GroupBy(group)
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

// UpdateParallelByCondition executes UPDATE query to all of the shards with conditions
func (xpr *XormParallel) UpdateParallelByCondition(objPtr interface{}, cond UpdateCondition) (int64, error) {
	// create session with the condition
	sessions := xpr.CreateUpdateSessions(cond)
	length := len(sessions)

	// execute query
	var errList []error
	results := make(chan int64, length)
	for _, s := range sessions {
		go func(s Session, obj interface{}) {
			count, err := s.Update(obj)
			if err != nil {
				errList = append(errList, err)
			}
			results <- count
		}(s, objPtr)
	}

	// wait for the results
	var counts int64
	for i := 0; i < length; i++ {
		v := <-results
		counts += v
	}
	if len(errList) > 0 {
		return counts, errors.NewErrParallelQuery(errList)
	}

	return counts, nil
}

// CreateUpdateSessions creates new sessions with conditional clause for UPDATE query
func (xpr *XormParallel) CreateUpdateSessions(cond UpdateCondition) []Session {
	var sessions []Session
	masters := xpr.orm.Masters(cond.Table)
	for _, master := range masters {
		s := master.NewSession()
		for _, w := range cond.Where {
			s.And(w.Statement, w.Args...)
		}
		for _, in := range cond.WhereIn {
			s.In(in.Statement, in.Args...)
		}

		if cond.AllColumns {
			s.AllCols()
		}
		for _, col := range cond.Columns {
			s.Cols(col)
		}
		for _, col := range cond.MustColumns {
			s.MustCols(col)
		}
		for _, col := range cond.OmitColumns {
			s.Omit(col)
		}
		for _, col := range cond.NullableColumns {
			s.Nullable(col)
		}

		for _, exp := range cond.Increments {
			s.Incr(exp.Statement, exp.Args...)
		}
		for _, exp := range cond.Decrements {
			s.Decr(exp.Statement, exp.Args...)
		}

		sessions = append(sessions, s)
	}
	return sessions
}
