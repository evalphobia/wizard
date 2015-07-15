package wizard

// UseMaster returns db master
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

// UseMasters returns all db master instances for sharding
func (w *Wizard) UseMasters(obj interface{}) []interface{} {
	c := w.getCluster(obj)
	if c == nil {
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

// UseSlave randomly returns db slave from the slaves
// if any slave is not set, master is returned
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

func (w *Wizard) UseSlaves(obj interface{}) []interface{} {
	c := w.getCluster(obj)
	if c == nil {
		return nil
	}
	var results []interface{}
	for _, node := range c.Slaves() {
		db := node.DB()
		if db == nil {
			continue
		}
		results = append(results, db)
	}
	return results
}

// UseMasterBySlot returns db master for sharding by hash slot id
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

// UseSlaveBySlot randomly returns db slave for sharding by hash slot id
// if any slave is not set in the cluster, master is returned
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
