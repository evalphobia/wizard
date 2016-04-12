package xorm

import "github.com/evalphobia/wizard"

// Xorm manages database sessions for xorm
type Xorm struct {
	*XormWizard
	*XormFunction
	*XormSessionManager
	*XormParallel

	Wiz *wizard.Wizard
}

// New creates initialized *Xorm
func New(wiz *wizard.Wizard) *Xorm {
	orm := &Xorm{}
	orm.Wiz = wiz
	orm.XormFunction = &XormFunction{orm: orm}
	orm.XormWizard = &XormWizard{wiz}
	orm.XormSessionManager = &XormSessionManager{
		orm:  orm,
		list: make(map[Identifier]*SessionList),
	}
	orm.XormParallel = &XormParallel{orm: orm}
	return orm
}
