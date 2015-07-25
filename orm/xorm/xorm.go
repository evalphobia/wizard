package xorm

import (
	"github.com/evalphobia/wizard"
)

// Xorm manages database sessions for xorm
type Xorm struct {
	*XormWizard
	*XormFunction
	*XormSession
	*XormTransaction
	*XormParallel

	Wiz      *wizard.Wizard
	readOnly bool
	autoTx   bool
}

// New creates initialized *Xorm
func New(wiz *wizard.Wizard) *Xorm {
	orm := &Xorm{}
	orm.Wiz = wiz
	orm.XormFunction = &XormFunction{orm: orm}
	orm.XormWizard = &XormWizard{wiz}
	orm.XormSession = &XormSession{
		orm:      orm,
		sessions: make(map[interface{}]Session),
	}
	orm.XormTransaction = &XormTransaction{
		orm:          orm,
		transactions: make(map[interface{}]Session),
	}
	orm.XormParallel = &XormParallel{orm: orm}
	return orm
}

// ReadOnly set write proof flag
func (orm *Xorm) ReadOnly(b bool) {
	orm.readOnly = b
}

// IsReadOnly checks in write proof mode or not
func (orm *Xorm) IsReadOnly() bool {
	return orm.readOnly
}

// SetAutoTransaction sets auto transaction flag
func (orm *Xorm) SetAutoTransaction(b bool) {
	orm.autoTx = b
}

// IsAutoTransaction checks in auto transaction mode or not
func (orm *Xorm) IsAutoTransaction() bool {
	return orm.autoTx
}
