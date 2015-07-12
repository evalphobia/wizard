package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

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

func (x *Xorm) Commit(obj interface{}) error {
	s, err := x.GetSession(obj)
	if err != nil {
		return err
	}
	x.DeleteSession(obj)
	return s.Commit()
}

func (x *Xorm) Rollback(obj interface{}) error {
	s, err := x.GetSession(obj)
	if err != nil {
		return err
	}
	x.DeleteSession(obj)
	return s.Rollback()
}

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

func (x *Xorm) getSession(db interface{}) Session {
	return x.sessions[db]
}

func (x *Xorm) addSession(db interface{}, s Session) {
	x.sessions[db] = s
}

func (x *Xorm) DeleteSession(obj interface{}) {
	db := x.UseMaster(obj)
	if db == nil {
		return
	}
	x.deleteSession(db)
}

func (x *Xorm) deleteSession(db interface{}) {
	delete(x.sessions, db)
}

func (x *Xorm) LazyBegin(obj interface{}) {
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return
	}
	x.lazy.prepareTx(NormalizeValue(obj))
}

func (x *Xorm) InLazyTx(obj interface{}) bool {
	return x.lazy.inLazyTx(NormalizeValue(obj))
}

func (x *Xorm) LazyBeginOne(obj interface{}) (Session, error) {
	n := NormalizeValue(obj)
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return nil, errors.NewErrNilDB(n)
	}
	db := x.UseMaster(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(n)
	}
	return x.lazy.BeginOne(n, db)
}

func (x *Xorm) LazyCommit(obj interface{}) error {
	n := NormalizeValue(obj)
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return errors.NewErrNilDB(n)
	}
	return x.lazy.CommitAll(n)
}

func (x *Xorm) LazyAbort(obj interface{}) error {
	n := NormalizeValue(obj)
	dbs := x.UseMasters(obj)
	if len(dbs) == 0 {
		return errors.NewErrNilDB(n)
	}
	return x.lazy.RollbackAll(n)
}
