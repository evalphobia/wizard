package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

// BeginSession returns session with transaction for the db of given object
func (x *Xorm) BeginSession(obj interface{}) (Session, error) {
	db := x.UseMaster(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	s := db.NewSession()
	err := s.Begin()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Begin starts the transaction for the db of given object
func (x *Xorm) Begin(obj interface{}) error {
	db := x.UseMaster(obj)
	if db == nil {
		return errors.NewErrNilDB(NormalizeValue(obj))
	}
	if x.getSession(db) != nil {
		return errors.NewErrDuplicateTx()
	}
	s := db.NewSession()
	err := s.Begin()
	if err != nil {
		return err
	}
	x.addSession(db, s)
	return nil
}

// Commit commits the transaction for the db of given object
func (x *Xorm) Commit(obj interface{}) error {
	s, err := x.GetSession(obj)
	if err != nil {
		return err
	}
	x.DeleteSession(obj)
	if x.readOnly {
		return nil
	}
	return s.Commit()
}

// Rollback aborts the transaction for the db of given object
func (x *Xorm) Rollback(obj interface{}) error {
	s, err := x.GetSession(obj)
	if err != nil {
		return err
	}
	x.DeleteSession(obj)
	return s.Rollback()
}

// GetOrCreateSession returns the session for the db of given object
// if no session exists for the object, create new one and return it
func (x *Xorm) GetOrCreateSession(obj interface{}) (Session, error) {
	if x.InLazyTx(obj) {
		return x.LazyBeginOne(obj)
	}
	db := x.UseMaster(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	s := x.getSession(db)
	if s == nil {
		return db.NewSession(), nil
	}
	return s, nil
}

// GetSession returns the session for the db of given object
func (x *Xorm) GetSession(obj interface{}) (Session, error) {
	db := x.UseMaster(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	s := x.getSession(db)
	if s == nil {
		return nil, errors.NewErrNoSession(NormalizeValue(obj))
	}
	return s, nil
}

func (x *Xorm) GetSessionByShardKey(obj interface{}, id int64) (Session, error) {
	db := x.UseMasterByShardKey(obj, id)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return db.NewSession(), nil
}

// getSession returns the session for the db
func (x *Xorm) getSession(db interface{}) Session {
	return x.sessions[db]
}

// addSession saves the session for the db
func (x *Xorm) addSession(db interface{}, s Session) {
	x.sessions[db] = s
}

// DeleteSession removes the saved session for the db of given object
func (x *Xorm) DeleteSession(obj interface{}) {
	db := x.UseMaster(obj)
	if db == nil {
		return
	}
	x.deleteSession(db)
}

// deleteSession removes the saved session for the db
func (x *Xorm) deleteSession(db interface{}) {
	delete(x.sessions, db)
}

// GetSlaveSession returns the session for the slave db of given object
func (x *Xorm) GetSlaveSession(obj interface{}) (Session, error) {
	db := x.UseSlave(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return db.NewSession(), nil
}

func (x *Xorm) GetSlaveSessionByShardKey(obj interface{}, id int64) (Session, error) {
	db := x.UseSlaveByShardKey(obj, id)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return db.NewSession(), nil
}

// LazyBegin set a transaction flag for the db of given object
func (x *Xorm) LazyBegin(obj interface{}) {
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return
	}
	x.lazy.prepareTx(NormalizeValue(obj))
}

// InLazyTx checks if a transaction flag in the db is on or not
func (x *Xorm) InLazyTx(obj interface{}) bool {
	return x.lazy.inLazyTx(NormalizeValue(obj))
}

// LazyBeginOne starts transaction for the db of given object
func (x *Xorm) LazyBeginOne(obj interface{}) (Session, error) {
	n := NormalizeValue(obj)
	db := x.UseMaster(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(n)
	}
	return x.lazy.BeginOne(n, db)
}

// LazyCommit commits all transaction for the db of given object
func (x *Xorm) LazyCommit(obj interface{}) error {
	n := NormalizeValue(obj)
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return errors.NewErrNilDB(n)
	}
	if x.readOnly {
		return nil
	}
	return x.lazy.CommitAll(n)
}

// LazyRollback aborts all transaction for the db of given object
func (x *Xorm) LazyRollback(obj interface{}) error {
	n := NormalizeValue(obj)
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return errors.NewErrNilDB(n)
	}
	return x.lazy.RollbackAll(n)
}
