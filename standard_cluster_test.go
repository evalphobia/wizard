package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCluster(t *testing.T) {
	assert := assert.New(t)

	var c *StandardCluster
	c = NewCluster("db")
	assert.IsType(Node{}, *c.master)
	assert.Equal("db", c.master.db, "db should be saved on NewCluster()")
}

func TestStandardClusterMaster(t *testing.T) {
	assert := assert.New(t)

	var c *StandardCluster
	var node *Node

	c = NewCluster("db")
	node = c.Master()
	assert.Equal("db", node.db)
	assert.Equal(c.master, node, "Master() shoud equal to StandardCluster.master")

	c.RegisterMaster("db2")
	node = c.Master()
	assert.Equal("db2", node.db)
}

func TestStandardClusterMasters(t *testing.T) {
	assert := assert.New(t)

	var c *StandardCluster
	var nodes []*Node

	c = NewCluster("db")
	nodes = c.Masters()
	assert.Equal(c.master, nodes[0])
	assert.Len(nodes, 1)
}

func TestStandardClusterSlave(t *testing.T) {
	assert := assert.New(t)

	var c *StandardCluster
	var node *Node

	c = NewCluster("master")
	node = c.Slave()
	assert.IsType(Node{}, *node)
	assert.Equal("master", node.db)
	assert.Equal(c.master, node, "Slave() shoud equal to StandardCluster.master when no slaves")
	assert.Len(c.slaves, 0)

	c.RegisterSlave("slave")
	node = c.Slave()
	assert.IsType(Node{}, *node)
	assert.Equal("slave", node.db)
	assert.Equal(c.slaves[0], node, "Slave() shoud equal to node in StandardCluster.slaves")
	assert.Len(c.slaves, 1)

	for i, max := 0, 100; i < max; i++ {
		c.RegisterSlave(i)
	}
	assert.Len(c.slaves, 101)

	node = c.Slave()
	db := node.db
	for i, max := 0, 10; i < max; i++ {
		node = c.Slave()
		if node.db != db {
			return
		}
	}
	t.Error("Slave() should return different nodes")
}

func TestStandardClusterSelectByKey(t *testing.T) {
	assert := assert.New(t)

	var c, c2, c3, c4 *StandardCluster

	c = NewCluster("db")
	c2 = c.SelectByKey(0)
	c3 = c.SelectByKey(1)
	c4 = c.SelectByKey(9999)
	assert.Equal(c, c2)
	assert.Equal(c, c3)
	assert.Equal(c, c4)
}

func TestStandardClusterRegisterMaster(t *testing.T) {
	assert := assert.New(t)

	var c *StandardCluster

	c = NewCluster("db")
	c.RegisterMaster("db2")
	c.RegisterMaster("db3")

	assert.Equal("db3", c.master.db)

	c.RegisterMaster("db4")
	assert.Equal("db4", c.master.db)
}

func TestStandardClusterRegisterSlave(t *testing.T) {
	assert := assert.New(t)

	var c *StandardCluster

	c = NewCluster("db")
	c.RegisterSlave("db2")
	c.RegisterSlave("db3")
	c.RegisterSlave("db4")

	assert.Equal("db2", c.slaves[0].db)
	assert.Equal("db3", c.slaves[1].db)
	assert.Equal("db4", c.slaves[2].db)
}
