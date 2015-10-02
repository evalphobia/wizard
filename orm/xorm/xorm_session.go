package xorm

import (
	"github.com/evalphobia/wizard/errors"
)

// XormSession manages database sessions for xorm
type XormSession struct {
	orm      *Xorm
	sessions map[interface{}]Session
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
	return xse.session(db, obj)
}

// NewMasterSessionByKey returns new master session by shard key
func (xse *XormSession) NewMasterSessionByKey(obj interface{}, key interface{}) (Session, error) {
	if xse.orm.IsAutoTransaction() {
		return xse.orm.TransactionByKey(obj, key)
	}
	db := xse.orm.MasterByKey(obj, key)
	return xse.session(db, obj)
}

// NewAllMasterSessions returns all of master sessions for the db of given object
func (xse *XormSession) NewAllMasterSessions(obj interface{}) ([]Session, error) {
	dbs := xse.orm.Masters(obj)

	var sessions []Session
	var errList []error
	for _, db := range dbs {
		var s Session
		var err error

		switch {
		case xse.orm.IsAutoTransaction():
			s, err = xse.orm.transaction(db, obj)
		default:
			s, err = xse.session(db, obj)
		}

		if err != nil {
			errList = append(errList, err)
			continue
		}
		sessions = append(sessions, s)
	}

	if len(errList) > 0 {
		return sessions, errors.NewErrNilDBs(errList)
	}
	return sessions, nil
}

// NewSlaveSession returns new slave session for the slave db of given object
func (xse *XormSession) NewSlaveSession(obj interface{}) (Session, error) {
	db := xse.orm.Slave(obj)
	return xse.session(db, obj)
}

// NewSlaveSessionByKey returns new slave session by shard key
func (xse *XormSession) NewSlaveSessionByKey(obj interface{}, key interface{}) (Session, error) {
	db := xse.orm.SlaveByKey(obj, key)
	return xse.session(db, obj)
}

// session returns the session for the db of given object
// if old session exists for the object, return it,
// if no session exists for the object, create new one and return it
func (xse *XormSession) session(db Engine, obj interface{}) (Session, error) {
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	// use old session
	s := xse.getSession(db)
	if s != nil {
		return s, nil
	}

	// create new session
	s = db.NewSession()
	xse.addSession(db, s)
	return s, nil
}

// getSession returns the session for the db
func (xse *XormSession) getSession(db interface{}) Session {
	return xse.sessions[db]
}

// addSession saves the session for the db
func (xse *XormSession) addSession(db interface{}, s Session) {
	xse.sessions[db] = s
}

// CloseAll closes all of sessions and engines
func (xse *XormSession) CloseAll() {
	for db, s := range xse.sessions {
		s.Close()
		if e, ok := db.(Engine); ok {
			e.Close()
		}
	}
	xse.sessions = make(map[interface{}]Session)
}
