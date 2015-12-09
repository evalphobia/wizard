package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

// ForceNewTransaction returns the session with new transaction
func (xse *XormSessionManager) ForceNewTransaction(obj interface{}) (Session, error) {
	db := xse.orm.Master(obj)
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
func (xse *XormSessionManager) Transaction(id Identifier, obj interface{}) (Session, error) {
	db := xse.orm.Master(obj)
	return xse.transaction(id, obj, db)
}

// TransactionByKey returns the session with transaction by shard key
func (xse *XormSessionManager) TransactionByKey(id Identifier, obj interface{}, key interface{}) (Session, error) {
	db := xse.orm.MasterByKey(obj, key)
	return xse.transaction(id, obj, db)
}

// transaction returns the session with transaction for the db of given object
// if old transaction exists for the object, return it,
// if no transaction exists for the object, create new one and return it
func (xse *XormSessionManager) transaction(id Identifier, obj interface{}, db Engine) (Session, error) {
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	// use old transaction
	s := xse.getTransactionFromList(id, db)
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
	xse.addTransactionIntoList(id, db, s)
	return s, nil
}

// getTransactionFromList returns the session with transaction for the db
func (xse *XormSessionManager) getTransactionFromList(id Identifier, db interface{}) Session {
	if !xse.hasSessionList(id) {
		return nil
	}
	sl := xse.getOrCreateSessionList(id)
	return sl.getTransaction(db)
}

// addTransactionIntoList saves the session with transaction for the db
func (xse *XormSessionManager) addTransactionIntoList(id Identifier, db interface{}, s Session) {
	xse.lock.Lock()
	defer xse.lock.Unlock()
	sl := xse.getOrCreateSessionList(id)
	sl.addTransaction(db, s)
}

// AutoTransaction starts transaction for the session and store it
// if not in the AutoTransaction mode, nothing happens
// if old transaction exists, return it
func (xse *XormSessionManager) AutoTransaction(id Identifier, obj interface{}, s Session) error {
	sl := xse.getOrCreateSessionList(id)
	if !sl.IsAutoTransaction() {
		return nil
	}
	db := xse.orm.Master(obj)
	oldTx := sl.getTransaction(db)
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

	sl.addTransaction(db, s)
	return nil
}

// CommitAll commits all of transactions
func (xse *XormSessionManager) CommitAll(id Identifier) error {
	xse.lock.Lock()
	defer xse.lock.Unlock()
	if !xse.hasSessionList(id) {
		return nil
	}

	sl := xse.getOrCreateSessionList(id)
	switch {
	case sl == nil:
		return nil
	case sl.IsReadOnly():
		return nil
	}

	var errList []error
	for _, s := range sl.getTransactions() {
		err := s.Commit()
		if err != nil {
			errList = append(errList, err)
		}
		s.Init()
	}

	sl.clearTransactions()
	if len(errList) > 0 {
		return errors.NewErrCommitAll(errList)
	}
	return nil
}

// RollbackAll aborts all of transactions
func (xse *XormSessionManager) RollbackAll(id Identifier) error {
	xse.lock.Lock()
	defer xse.lock.Unlock()
	if !xse.hasSessionList(id) {
		return nil
	}

	sl := xse.getOrCreateSessionList(id)
	switch {
	case sl == nil:
		return nil
	case sl.IsReadOnly():
		return nil
	}

	var errList []error
	for _, s := range sl.getTransactions() {
		err := s.Rollback()
		if err != nil {
			errList = append(errList, err)
		}
		s.Init()
	}

	sl.clearTransactions()
	if len(errList) > 0 {
		return errors.NewErrRollbackAll(errList)
	}
	return nil
}
