package xorm

import (
// "github.com/evalphobia/wizard/errors"
)

type SessionList struct {
	readOnly bool
	autoTx   bool

	sessions     map[interface{}]Session
	transactions map[interface{}]Session
}

func newSessionList() *SessionList {
	return &SessionList{
		sessions:     make(map[interface{}]Session),
		transactions: make(map[interface{}]Session),
	}
}

func (l *SessionList) hasSession(db interface{}) bool {
	_, ok := l.sessions[db]
	return ok
}

func (l *SessionList) getSession(db interface{}) Session {
	return l.sessions[db]
}

func (l *SessionList) addSession(db interface{}, s Session) {
	l.sessions[db] = s
}

func (l *SessionList) getSessions() map[interface{}]Session {
	return l.sessions
}

func (l *SessionList) clearSessions() {
	l.sessions = make(map[interface{}]Session)
}

func (l *SessionList) getTransaction(db interface{}) Session {
	return l.transactions[db]
}

func (l *SessionList) addTransaction(db interface{}, s Session) {
	l.transactions[db] = s
}

func (l *SessionList) getTransactions() map[interface{}]Session {
	return l.transactions
}

func (l *SessionList) clearTransactions() {
	l.transactions = make(map[interface{}]Session)
}

// ReadOnly set write proof flag
func (l *SessionList) ReadOnly(b bool) {
	l.readOnly = b
}

// IsReadOnly checks in write proof mode or not
func (l *SessionList) IsReadOnly() bool {
	return l.readOnly
}

// SetAutoTransaction sets auto transaction flag
func (l *SessionList) SetAutoTransaction(b bool) {
	l.autoTx = b
}

// IsAutoTransaction checks in auto transaction mode or not
func (l *SessionList) IsAutoTransaction() bool {
	return l.autoTx
}
