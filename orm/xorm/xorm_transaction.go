package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

// XormTransaction manages transaction of session for xorm
type XormTransaction struct {
	orm          *Xorm
	transactions map[interface{}]Session
}

// ForceNewTransaction returns the session with new transaction
func (xtx *XormTransaction) ForceNewTransaction(obj interface{}) (Session, error) {
	db := xtx.orm.Master(obj)
	s, err := newSession(db, obj)
	if err != nil {
		return nil, err
	}
	err = s.Begin()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Transaction returns the session with transaction for the db of given object
func (xtx *XormTransaction) Transaction(obj interface{}) (Session, error) {
	db := xtx.orm.Master(obj)
	return xtx.transaction(db, obj)
}

// TransactionByKey returns the session with transaction by shard key
func (xtx *XormTransaction) TransactionByKey(obj interface{}, key interface{}) (Session, error) {
	db := xtx.orm.MasterByKey(obj, key)
	return xtx.transaction(db, obj)
}

// transaction returns the session with transaction for the db of given object
// if old transaction exists for the object, return it,
// if no transaction exists for the object, create new one and return it
func (xtx *XormTransaction) transaction(db Engine, obj interface{}) (Session, error) {
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	// use old transaction
	s := xtx.getTransaction(db)
	if s != nil {
		return s, nil
	}

	// create new transaction
	s = db.NewSession()
	err := s.Begin()
	if err != nil {
		return nil, err
	}

	// save created session with transaction
	xtx.addTransaction(db, s)
	return s, nil
}

// getTransaction returns the session with transaction for the db
func (xtx *XormTransaction) getTransaction(db interface{}) Session {
	return xtx.transactions[db]
}

// addTransaction saves the session with transaction for the db
func (xtx *XormTransaction) addTransaction(db interface{}, s Session) {
	xtx.transactions[db] = s
}

// AutoTransaction starts transaction for the session and store it
// if not in the AutoTransaction mode, nothing happens
// if old transaction exits, return it
func (xtx *XormTransaction) AutoTransaction(obj interface{}, s Session) error {
	if !xtx.orm.IsAutoTransaction() {
		return nil
	}

	db := xtx.orm.Master(obj)
	oldTx := xtx.getTransaction(db)
	switch {
	case oldTx == s:
		return nil
	case oldTx != nil:
		return errors.NewErrAnotherTx(NormalizeValue(obj))
	}

	err := s.Begin()
	if err != nil {
		return err
	}

	xtx.addTransaction(db, s)
	return nil
}

// CommitAll commits all of transactions
func (xtx *XormTransaction) CommitAll() error {
	if xtx.orm.IsReadOnly() {
		return nil
	}
	var errList []error
	for _, s := range xtx.transactions {
		err := s.Commit()
		if err != nil {
			errList = append(errList, err)
		}
	}

	xtx.transactions = make(map[interface{}]Session)
	if len(errList) > 0 {
		return errors.NewErrCommitAll(errList)
	}
	return nil
}

// RollbackAll aborts all of transactions
func (xtx *XormTransaction) RollbackAll() error {
	if xtx.orm.IsReadOnly() {
		return nil
	}
	var errList []error
	for _, s := range xtx.transactions {
		err := s.Rollback()
		if err != nil {
			errList = append(errList, err)
		}
	}

	xtx.transactions = make(map[interface{}]Session)
	if len(errList) > 0 {
		return errors.NewErrRollbackAll(errList)
	}
	return nil
}
