package wizard

import (
	"github.com/evalphobia/wizard/errors"
)

// Cluster is interface for [StandardCluster | ShardCluster]
type Cluster interface {
	SelectByKey(interface{}) *StandardCluster
	Master() *Node
	Masters() []*Node
	Slave() *Node
	Slaves() []*Node
}

// Wizard manages all the database cluster for your app
type Wizard struct {
	clusters       map[interface{}]Cluster
	defaultCluster Cluster
}

// NewWizard returns initialized empty Wizard
func NewWizard() *Wizard {
	return &Wizard{
		clusters: make(map[interface{}]Cluster),
	}
}

// SetDefault set default cluster
// if default is set, this cluster acts like catchall, handles all the other tables.
func (w *Wizard) SetDefault(c Cluster) {
	w.defaultCluster = c
}

// HasDefault checks default cluster is set or not
func (w *Wizard) HasDefault() bool {
	return w.defaultCluster != nil
}

// getCluster returns the cluster by name mapping
func (w *Wizard) getCluster(obj interface{}) Cluster {
	c, ok := w.clusters[NormalizeValue(obj)]
	switch {
	case ok:
		return c
	case w.HasDefault():
		return w.defaultCluster
	default:
		return nil
	}
}

// RegisterTables adds cluster and tables for name mapping
func (w *Wizard) RegisterTables(c Cluster, list ...interface{}) error {
	for _, obj := range list {
		v := NormalizeValue(obj)
		if _, ok := w.clusters[v]; ok {
			return errors.NewErrAlreadyRegistared(v)
		}
		w.clusters[v] = c
	}
	return nil
}

// setCluster set the cluster with name mapping
func (w *Wizard) setCluster(c Cluster, obj interface{}) {
	w.clusters[NormalizeValue(obj)] = c
}

// CreateCluster set and returns the new StandardCluster
func (w *Wizard) CreateCluster(obj interface{}, db interface{}) *StandardCluster {
	c := NewCluster(db)
	w.setCluster(c, obj)
	return c
}

// CreateShardCluster set and returns the new ShardCluster
func (w *Wizard) CreateShardCluster(obj interface{}, slot int64) *ShardCluster {
	if slot < 1 {
		slot = 1
	}
	c := &ShardCluster{
		slotsize: slot,
	}
	w.setCluster(c, obj)
	return c
}

// Select returns StandardCluster by name mapping (and implicit hash slot from struct field)
func (w *Wizard) Select(obj interface{}) *StandardCluster {
	c := w.getCluster(obj)
	switch v := c.(type) {
	case *StandardCluster:
		return v
	case *ShardCluster:
		return v.SelectByKey(getShardKey(obj))
	default:
		return nil
	}
}

// SelectByKey returns StandardCluster by name mapping and shard key
func (w *Wizard) SelectByKey(obj interface{}, key interface{}) *StandardCluster {
	c := w.getCluster(obj)
	if c == nil {
		return nil
	}
	return c.SelectByKey(key)
}
