package xorm

import (
	"github.com/evalphobia/wizard"
	"github.com/evalphobia/wizard/errors"
)

// Xorm manages database sessions for xorm
type Xorm struct {
	c        *wizard.Wizard
	sessions map[interface{}]Session
	lazy     *LazySessionList
	readOnly bool
}

// New creates initialized *Xorm
func New(c *wizard.Wizard) *Xorm {
	return &Xorm{
		c:        c,
		sessions: make(map[interface{}]Session),
		lazy:     newLazySessionList(),
	}
}

func (x *Xorm) ReadOnly(b bool) {
	x.readOnly = b
}

// UseMaster returns master db for the given object
func (x *Xorm) UseMaster(obj interface{}) Engine {
	db := x.c.UseMaster(obj)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

func (x *Xorm) UseMasterByShardKey(obj interface{}, id int64) Engine {
	db := x.c.UseMasterBySlot(obj, id)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

// UseMasters returns all of sharded master db for the given object
func (x *Xorm) UseMasters(obj interface{}) []Engine {
	var results []Engine
	for _, db := range x.c.UseMasters(obj) {
		e, ok := db.(Engine)
		if !ok {
			continue
		}
		results = append(results, e)
	}
	return results
}

// UseSlave randomly returns one of the slave db for the given object
func (x *Xorm) UseSlave(obj interface{}) Engine {
	db := x.c.UseSlave(obj)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

func (x *Xorm) UseSlaveByShardKey(obj interface{}, id int64) Engine {
	db := x.c.UseSlaveBySlot(obj, id)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

func (x *Xorm) UseSlaves(obj interface{}) []Engine {
	var results []Engine
	for _, db := range x.c.UseSlaves(obj) {
		e, ok := db.(Engine)
		if !ok {
			continue
		}
		results = append(results, e)
	}
	return results
}

// Get executes xorm.Sessions.Get() in slave db
func (x *Xorm) Get(obj interface{}, fn func(Session) (bool, error)) (bool, error) {
	db := x.UseSlave(obj)
	if db == nil {
		return false, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return fn(db.NewSession())
}

// Find executes xorm.Sessions.Find() in slave db
func (x *Xorm) Find(obj interface{}, fn func(Session) error) error {
	db := x.UseSlave(obj)
	if db == nil {
		return errors.NewErrNilDB(NormalizeValue(obj))
	}
	return fn(db.NewSession())
}

// Count executes xorm.Sessions.Count() in slave db
func (x *Xorm) Count(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	db := x.UseSlave(obj)
	if db == nil {
		return 0, errors.NewErrNilDB(NormalizeValue(obj))
	}
	return fn(db.NewSession())
}

// Insert executes xorm.Sessions.Insert() in master db
func (x *Xorm) Insert(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	if x.readOnly {
		return 0, nil
	}

	s, err := x.GetOrCreateSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}

// Update executes xorm.Sessions.Update() in master db
func (x *Xorm) Update(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	if x.readOnly {
		return 0, nil
	}

	s, err := x.GetOrCreateSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}

// GetUsingMaster executes xorm.Sessions.Get() in master db
func (x *Xorm) GetUsingMaster(obj interface{}, fn func(Session) (bool, error)) (bool, error) {
	s, err := x.GetOrCreateSession(obj)
	if err != nil {
		return false, err
	}
	return fn(s)
}

// FindUsingMaster executes xorm.Sessions.Find() in master db
func (x *Xorm) FindUsingMaster(obj interface{}, fn func(Session) error) error {
	s, err := x.GetOrCreateSession(obj)
	if err != nil {
		return err
	}
	return fn(s)
}

// CountUsingMaster executes xorm.Sessions.Count() in master db
func (x *Xorm) CountUsingMaster(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
	s, err := x.GetOrCreateSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}
