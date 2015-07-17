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
	vt := reflect.TypeOf(listPtr)
	if vt.Kind() != reflect.Ptr {
		return errors.NewErrArgType("listPtr must be a pointer")
	}
	elem := vt.Elem()
	if elem.Kind() != reflect.Slice && elem.Kind() != reflect.Map {
		return errors.NewErrArgType("listPtr must be a pointer of slice or map")
	}

	slaves := xpr.orm.Slaves(table)

	results := make(chan reflect.Value, len(slaves))
	var errList []error
	for _, slave := range slaves {
		s := slave.NewSession()
		list := reflect.New(elem)
		go func(s Session, list reflect.Value) {
			s.And(where, args...)
			err := s.Find(list.Interface())
			if err != nil {
				errList = append(errList, err)
			}
			results <- list
		}(s, list)
	}

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
