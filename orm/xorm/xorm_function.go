package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

// XormFunction manages xorm functions
type XormFunction struct {
	orm *Xorm
}

// Get executes xorm.Sessions.Get() in slave db
func (xfn XormFunction) Get(obj interface{}, fn func(Session) (bool, error)) (bool, error) {
	db := xfn.orm.Slave(obj)
	if db == nil {
		return false, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return fn(db.NewSession())
}

// Find executes xorm.Sessions.Find() in slave db
func (xfn XormFunction) Find(obj interface{}, fn func(Session) error) error {
	db := xfn.orm.Slave(obj)
	if db == nil {
		return errors.NewErrNilDB(NormalizeValue(obj))
	}
	return fn(db.NewSession())
}

// Count executes xorm.Sessions.Count() in slave db
func (xfn XormFunction) Count(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	db := xfn.orm.Slave(obj)
	if db == nil {
		return 0, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return fn(db.NewSession())
}

// Insert executes xorm.Sessions.Insert() in master db
func (xfn XormFunction) Insert(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	if xfn.orm.IsReadOnly() {
		return 0, nil
	}

	s, err := xfn.orm.NewMasterSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}

// Update executes xorm.Sessions.Update() in master db
func (xfn XormFunction) Update(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	if xfn.orm.IsReadOnly() {
		return 0, nil
	}

	s, err := xfn.orm.NewMasterSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}

// GetUsingMaster executes xorm.Sessions.Get() in master db
func (xfn XormFunction) GetUsingMaster(obj interface{}, fn func(Session) (bool, error)) (bool, error) {
	s, err := xfn.orm.NewMasterSession(obj)
	if err != nil {
		return false, err
	}
	return fn(s)
}

// FindUsingMaster executes xorm.Sessions.Find() in master db
func (xfn XormFunction) FindUsingMaster(obj interface{}, fn func(Session) error) error {
	s, err := xfn.orm.NewMasterSession(obj)
	if err != nil {
		return err
	}
	return fn(s)
}

// CountUsingMaster executes xorm.Sessions.Count() in master db
func (xfn XormFunction) CountUsingMaster(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	s, err := xfn.orm.NewMasterSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}
