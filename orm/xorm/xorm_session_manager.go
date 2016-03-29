package xorm

import (
	"sync"

	"github.com/evalphobia/wizard/errors"
)

// XormSessionManager manages database session list for xorm
type XormSessionManager struct {
	orm  *Xorm
	lock sync.RWMutex
	list map[Identifier]*SessionList
}

// Identifier is unique object for using same sessions
// e.g. *http.Request, context.Context, etc...
type Identifier interface{}

func newSession(db Engine, obj interface{}) (Session, error) {
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return db.NewSession(), nil
}

func (xse *XormSessionManager) SetAutoTransaction(id Identifier, b bool) {
	sl := xse.getOrCreateSessionList(id)
	sl.SetAutoTransaction(b)
}

func (xse *XormSessionManager) IsAutoTransaction(id Identifier) bool {
	sl := xse.getOrCreateSessionList(id)
	return sl.IsAutoTransaction()
}

func (xse *XormSessionManager) ReadOnly(id Identifier, b bool) {
	sl := xse.getOrCreateSessionList(id)
	sl.ReadOnly(b)
}

func (xse *XormSessionManager) IsReadOnly(id Identifier) bool {
	sl := xse.getOrCreateSessionList(id)
	return sl.IsReadOnly()
}

// NewMasterSession returns new master session for the db of given object
func (xse *XormSessionManager) NewMasterSession(obj interface{}) (Session, error) {
	db := xse.orm.Master(obj)
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}

	return db.NewSession(), nil
}

// NewMasterSession returns new master session for the db of given object
func (xse *XormSessionManager) UseMasterSession(id Identifier, obj interface{}) (Session, error) {
	db := xse.orm.Master(obj)
	sl := xse.getOrCreateSessionList(id)
	if sl.IsAutoTransaction() {
		return xse.transaction(id, obj, db)
	}
	return xse.session(id, obj, db)
}

// UseMasterSessionByKey returns new master session by shard key
func (xse *XormSessionManager) UseMasterSessionByKey(id Identifier, obj interface{}, key interface{}) (Session, error) {
	db := xse.orm.MasterByKey(obj, key)
	sl := xse.getOrCreateSessionList(id)
	if sl.IsAutoTransaction() {
		return xse.transaction(id, obj, db)
	}
	return xse.session(id, obj, db)
}

// UseAllMasterSessions returns all of master sessions for the db of given object
func (xse *XormSessionManager) UseAllMasterSessions(id Identifier, obj interface{}) ([]Session, error) {
	dbs := xse.orm.Masters(obj)

	var sessions []Session
	var errList []error
	for _, db := range dbs {
		var s Session
		var err error

		switch {
		// case xse.orm.IsAutoTransaction():
		// s, err = xse.orm.transaction(obj, db)
		default:
			s, err = xse.session(id, obj, db)
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

// UseSlaveSession returns new slave session for the slave db of given object
func (xse *XormSessionManager) UseSlaveSession(id Identifier, obj interface{}) (Session, error) {
	db := xse.orm.Slave(obj)
	return xse.session(id, obj, db)
}

// UseSlaveSessionByKey returns new slave session by shard key
func (xse *XormSessionManager) UseSlaveSessionByKey(id Identifier, obj interface{}, key interface{}) (Session, error) {
	db := xse.orm.SlaveByKey(obj, key)
	return xse.session(id, obj, db)
}

// session returns the session for the db of given object
// if old session exists for the object, return it,
// if no session exists for the object, create new one and return it
func (xse *XormSessionManager) session(id Identifier, obj interface{}, db Engine) (Session, error) {
	if db == nil {
		return nil, errors.NewErrNilDB(NormalizeValue(obj))
	}
	// use old session

	s := xse.getSessionFromList(id, db)
	if s != nil {
		return s, nil
	}

	// create new session
	s = db.NewSession()
	xse.addSessionIntoList(id, db, s)
	return s, nil
}

// getSessionFromList returns the session for the db
func (xse *XormSessionManager) getSessionFromList(id Identifier, db interface{}) Session {
	if !xse.hasSessionList(id) {
		return nil
	}

	sl := xse.getOrCreateSessionList(id)
	return sl.getSession(db)
}

// addSessionIntoList saves the session for the db
func (xse *XormSessionManager) addSessionIntoList(id Identifier, db interface{}, s Session) {
	xse.lock.Lock()
	defer xse.lock.Unlock()
	sl := xse.getOrCreateSessionList(id)
	sl.addSession(db, s)
}

// CloseAll closes all of sessions and engines
func (xse *XormSessionManager) CloseAll(id Identifier) {
	xse.lock.Lock()
	defer xse.lock.Unlock()
	sl := xse.getOrCreateSessionList(id)

	for _, s := range sl.getSessions() {
		s.Close()
	}
	for _, s := range sl.getTransactions() {
		s.Close()
	}
	sl.clearSessions()
	sl.clearTransactions()
	delete(xse.list, id)
}

func (xse *XormSessionManager) newSessionList(id Identifier) *SessionList {
	if xse.list == nil {
		xse.list = make(map[Identifier]*SessionList)
	}
	xse.list[id] = newSessionList()
	return xse.list[id]
}

func (xse *XormSessionManager) hasSessionList(id Identifier) bool {
	_, ok := xse.list[id]
	return ok
}

func (xse *XormSessionManager) getOrCreateSessionList(id Identifier) *SessionList {
	if !xse.hasSessionList(id) {
		xse.newSessionList(id)
	}
	return xse.list[id]
}
