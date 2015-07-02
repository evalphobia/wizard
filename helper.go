package wizard

import (
	"github.com/evalphobia/go-log-wrapper/log"
)

var _ = log.Nothing

func (w *Wizard) UseMaster(obj interface{}) interface{} {
	cluster := w.Select(obj)
	if cluster == nil {
		return nil
	}
	db := cluster.Master()
	if db == nil {
		return nil
	}
	return db.DB()
}

func (w *Wizard) UseMasters(obj interface{}) []interface{} {
	c, ok := w.clusters[NormalizeValue(obj)]
	if !ok {
		return nil
	}
	var results []interface{}
	for _, node := range c.Masters() {
		db := node.DB()
		if db == nil {
			continue
		}
		results = append(results, db)
	}
	return results
}

func (w *Wizard) UseSlave(obj interface{}) interface{} {
	cluster := w.Select(obj)
	if cluster == nil {
		return nil
	}
	db := cluster.Slave()
	if db == nil {
		return nil
	}
	return db.DB()
}

func (w *Wizard) UseMasterBySlot(obj interface{}, id int64) interface{} {
	cluster := w.SelectBySlot(obj, id)
	if cluster == nil {
		return nil
	}
	db := cluster.Master()
	if db == nil {
		return nil
	}
	return db.DB()
}

func (w *Wizard) UseSlaveBySlot(obj interface{}, id int64) interface{} {
	cluster := w.SelectBySlot(obj, id)
	if cluster == nil {
		return nil
	}
	db := cluster.Slave()
	if db == nil {
		return nil
	}
	return db.DB()
}
