package wizard

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StandardCluster struct {
	table interface{}

	master *Node
	slaves []*Node
}

func NewCluster(db interface{}) *StandardCluster {
	node := NewNode(db)
	return &StandardCluster{master: node}
}

func (c StandardCluster) Table() interface{} {
	return c.table
}

func (c StandardCluster) Master() *Node {
	return c.master
}

func (c StandardCluster) Masters() []*Node {
	return []*Node{c.master}
}

func (c StandardCluster) Slave() *Node {
	if len(c.slaves) == 0 {
		return c.master
	}
	return c.slaves[rand.Intn(len(c.slaves))]
}

func (c *StandardCluster) SelectBySlot(i int64) *StandardCluster {
	return c
}

func (c *StandardCluster) RegisterMaster(db interface{}) {
	c.master = &Node{db: db}
}

func (c *StandardCluster) RegisterSlave(db interface{}) {
	c.slaves = append(c.slaves, &Node{db: db})
}
