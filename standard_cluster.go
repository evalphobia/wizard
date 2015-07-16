package wizard

import (
	"math/rand"
	"time"
)

func init() {
	// initialized for slave balancing
	rand.Seed(time.Now().UnixNano())
}

// StandardCluster is struct for typical(non-sharded) database cluster
type StandardCluster struct {
	master *Node
	slaves []*Node
}

// NewCluster returns the StandardCluster initialized with master database
func NewCluster(db interface{}) *StandardCluster {
	node := NewNode(db)
	return &StandardCluster{master: node}
}

// Master returns master database
func (c StandardCluster) Master() *Node {
	return c.master
}

// Masters is dummy method for interface
func (c StandardCluster) Masters() []*Node {
	return []*Node{c.master}
}

// Slave ramdomly returns the slave database.
// if no slave is registered, master is returned
func (c StandardCluster) Slave() *Node {
	if len(c.slaves) == 0 {
		return c.master
	}
	return c.slaves[rand.Intn(len(c.slaves))]
}

// Masters is dummy method for interface
func (c StandardCluster) Slaves() []*Node {
	return []*Node{c.Slave()}
}

// SelectByKey is dummy method for interface
func (c *StandardCluster) SelectByKey(v interface{}) *StandardCluster {
	return c
}

// RegisterMaster set new master node
func (c *StandardCluster) RegisterMaster(db interface{}) {
	c.master = &Node{db: db}
}

// RegisterSlave adds slave node
func (c *StandardCluster) RegisterSlave(db interface{}) {
	c.slaves = append(c.slaves, &Node{db: db})
}
