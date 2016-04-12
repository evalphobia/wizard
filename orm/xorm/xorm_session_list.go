package xorm

import "sync"

// SessionList contains db sessions list for one group
type SessionList struct {
	readOnly bool
	autoTx   bool

	sessMu   sync.RWMutex
	sessions map[interface{}]Session

	txMu         sync.RWMutex
	transactions map[interface{}]Session
}

func newSessionList() *SessionList {
	return &SessionList{
		sessions:     make(map[interface{}]Session),
		transactions: make(map[interface{}]Session),
	}
}

func (l *SessionList) hasSession(db interface{}) bool {
	l.sessMu.RLock()
	defer l.sessMu.RUnlock()
	_, ok := l.sessions[db]
	return ok
}

func (l *SessionList) getSession(db interface{}) Session {
	l.sessMu.RLock()
	defer l.sessMu.RUnlock()
	return l.sessions[db]
}

func (l *SessionList) addSession(db interface{}, s Session) {
	l.sessMu.Lock()
	defer l.sessMu.Unlock()
	l.sessions[db] = s
}

func (l *SessionList) getSessions() map[interface{}]Session {
	return l.sessions
}

func (l *SessionList) clearSessions() {
	l.sessMu.Lock()
	defer l.sessMu.Unlock()
	l.sessions = make(map[interface{}]Session)
}

func (l *SessionList) getTransaction(db interface{}) Session {
	l.txMu.RLock()
	defer l.txMu.RUnlock()
	return l.transactions[db]
}

func (l *SessionList) addTransaction(db interface{}, s Session) {
	l.txMu.Lock()
	defer l.txMu.Unlock()
	l.transactions[db] = s
}

func (l *SessionList) getTransactions() map[interface{}]Session {
	return l.transactions
}

func (l *SessionList) clearTransactions() {
	l.txMu.Lock()
	defer l.txMu.Unlock()
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
