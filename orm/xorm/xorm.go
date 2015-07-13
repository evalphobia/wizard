package xorm

import (
	"github.com/go-xorm/xorm"

	"github.com/evalphobia/wizard"
	"github.com/evalphobia/wizard/errors"
)

// Xorm manages database sessions for xorm
type Xorm struct {
	c        *wizard.Wizard
	sessions map[interface{}]Session
	lazy     *LazySessionList
}

// New creates initialized *Xorm
func New(c *wizard.Wizard) *Xorm {
	return &Xorm{
		c:        c,
		sessions: make(map[interface{}]Session),
		lazy:     newLazySessionList(),
	}
}

// UseMaster returns master db for the given object
func (x *Xorm) UseMaster(obj interface{}) *xorm.Engine {
	db := x.c.UseMaster(obj)
	if db == nil {
		return nil
	}
	return db.(*xorm.Engine)
}

// UseMasters returns all of sharded master db for the given object
func (x *Xorm) UseMasters(obj interface{}) []*xorm.Engine {
	var results []*xorm.Engine
	for _, db := range x.c.UseMasters(obj) {
		e, ok := db.(*xorm.Engine)
		if !ok {
			continue
		}
		results = append(results, e)
	}
	return results
}

// UseSlave randomly returns one of the slave db for the given object
func (x *Xorm) UseSlave(obj interface{}) *xorm.Engine {
	db := x.c.UseSlave(obj)
	if db == nil {
		return nil
	}
	return db.(*xorm.Engine)
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
	s, err := x.GetOrCreateSession(obj)
	if err != nil {
		return 0, err
	}
	return fn(s)
}

// Update executes xorm.Sessions.Update() in master db
func (x *Xorm) Update(obj interface{}, fn func(Session) (int64, error)) (int64, error) {
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
