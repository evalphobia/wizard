package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

// XormSession manages database sessions for xorm
type XormSession struct {
	orm *Xorm
}

func newSession(db Engine, obj interface{}) (Session, error) {
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return db.NewSession(), nil
}

// NewMasterSession returns new master session for the db of given object
func (xse *XormSession) NewMasterSession(obj interface{}) (Session, error) {
	if xse.orm.IsAutoTransaction() {
		return xse.orm.Transaction(obj)
	}
	db := xse.orm.Master(obj)
	return newSession(db, obj)
}

// NewMasterSessionByKey returns new master session by shard key
func (xse *XormSession) NewMasterSessionByKey(obj interface{}, key interface{}) (Session, error) {
	if xse.orm.IsAutoTransaction() {
		return xse.orm.TransactionByKey(obj, key)
	}
	db := xse.orm.MasterByKey(obj, key)
	return newSession(db, obj)
}

// NewSlaveSession returns new slave session for the slave db of given object
func (xse *XormSession) NewSlaveSession(obj interface{}) (Session, error) {
	db := xse.orm.Slave(obj)
	return newSession(db, obj)
}

// NewSlaveSessionByKey returns new slave session by shard key
func (xse *XormSession) NewSlaveSessionByKey(obj interface{}, key interface{}) (Session, error) {
	db := xse.orm.SlaveByKey(obj, key)
	return newSession(db, obj)
}
