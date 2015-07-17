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
	var results []interface{}
	c := w.getCluster(obj)
	if c == nil {
		return results
	}
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

// UseSlaves randomly returns all db slave instances for sharding
func (w *Wizard) UseSlaves(obj interface{}) []interface{} {
	var results []interface{}
	c := w.getCluster(obj)
	if c == nil {
		return results
	}
	for _, node := range c.Slaves() {
		db := node.DB()
		if db == nil {
			continue
		}
		results = append(results, db)
	}
	return results
}

// UseMasterByKey returns db master for sharding by shard key
func (w *Wizard) UseMasterByKey(obj interface{}, key interface{}) interface{} {
	cluster := w.SelectByKey(obj, key)
	if cluster == nil {
		return nil
	}
	db := cluster.Master()
	if db == nil {
		return nil
	}
	return db.DB()
}

// UseSlaveByKey randomly returns db slave for sharding by shard key
// if any slave is not set in the cluster, master is returned
func (w *Wizard) UseSlaveByKey(obj interface{}, key interface{}) interface{} {
	cluster := w.SelectByKey(obj, key)
	if cluster == nil {
		return nil
	}
	db := cluster.Slave()
	if db == nil {
		return nil
	}
	return db.DB()
}
