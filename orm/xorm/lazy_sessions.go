package xorm

import (
	"github.com/go-xorm/xorm"

	"github.com/evalphobia/wizard/errors"
)

// LazySessionList enables lazy transation
type LazySessionList struct {
	list map[interface{}]*LazySessions
	inTx map[interface{}]bool
}

func newLazySessionList() *LazySessionList {
	return &LazySessionList{
		list: make(map[interface{}]*LazySessions),
		inTx: make(map[interface{}]bool),
	}
}

func (l *LazySessionList) prepareTx(name interface{}) {
	l.list[name] = newLazySessions()
	l.inTx[name] = true
}

func (l *LazySessionList) inLazyTx(name interface{}) bool {
	b, ok := l.inTx[name]
	if !ok {
		return false
	}
	return b
}

func (l *LazySessionList) BeginOne(name interface{}, db *xorm.Engine) (Session, error) {
	ls := l.list[name]
	if ls == nil {
		return nil, errors.NewErrWrongTx()
	}
	return ls.getOrNewSession(db)
}

func (l *LazySessionList) CommitAll(name interface{}) error {
	ls := l.list[name]
	if ls == nil {
		return errors.NewErrWrongTx()
	}
	err := ls.CommitAll()
	l.inTx[name] = false
	l.list = make(map[interface{}]*LazySessions)
	return err
}

func (l *LazySessionList) RollbackAll(name interface{}) error {
	ls := l.list[name]
	if ls == nil {
		return errors.NewErrWrongTx()
	}
	err := ls.RollbackAll()
	l.inTx[name] = false
	l.list = make(map[interface{}]*LazySessions)
	return err
}

// LazySessions
type LazySessions struct {
	sessions map[interface{}]Session
}

func newLazySessions() *LazySessions {
	return &LazySessions{
		sessions: make(map[interface{}]Session),
	}
}

func (ls *LazySessions) getOrNewSession(db *xorm.Engine) (Session, error) {
	s := ls.sessions[db]
	if s != nil {
		return s, errors.NewErrDuplicateTx()
	}
	s = db.NewSession()
	err := s.Begin()
	if err != nil {
		return nil, err
	}
	ls.sessions[db] = s
	return s, nil
}

func (ls *LazySessions) CommitAll() error {
	var errs []error
	for _, s := range ls.sessions {
		err := s.Commit()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.NewErrCommitAll(errs)
}

func (ls *LazySessions) RollbackAll() error {
	var errs []error
	for _, s := range ls.sessions {
		err := s.Rollback()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.NewErrCommitAll(errs)
}
