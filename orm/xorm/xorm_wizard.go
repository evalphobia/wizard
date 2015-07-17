package xorm

import (
	"github.com/evalphobia/wizard"
)

// XormWizard is struct for database selector
type XormWizard struct {
	*wizard.Wizard
}

// Master returns master db for the given object
func (xwiz XormWizard) Master(obj interface{}) Engine {
	db := xwiz.UseMaster(obj)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

// MasterByKey returns master db by shard key
func (xwiz XormWizard) MasterByKey(obj interface{}, key interface{}) Engine {
	db := xwiz.UseMasterByKey(obj, key)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

// Masters returns all of sharded master db for the given object
func (xwiz XormWizard) Masters(obj interface{}) []Engine {
	var results []Engine
	for _, db := range xwiz.UseMasters(obj) {
		e, ok := db.(Engine)
		if !ok || e == nil {
			continue
		}
		results = append(results, e)
	}
	return results
}

// Slave randomly returns one of the slave db for the given object
func (xwiz XormWizard) Slave(obj interface{}) Engine {
	db := xwiz.UseSlave(obj)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

// SlaveByKey randomly returns one of the slave db by shard key
func (xwiz XormWizard) SlaveByKey(obj interface{}, key interface{}) Engine {
	db := xwiz.UseSlaveByKey(obj, key)
	if db == nil {
		return nil
	}
	return db.(Engine)
}

// Slaves randomly returns all of sharded slave db for the given object
func (xwiz XormWizard) Slaves(obj interface{}) []Engine {
	var results []Engine
	for _, db := range xwiz.UseSlaves(obj) {
		e, ok := db.(Engine)
		if !ok || e == nil {
			continue
		}
		results = append(results, e)
	}
	return results
}
