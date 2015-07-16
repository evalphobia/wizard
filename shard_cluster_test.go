package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCreateCluster(prefix string) *StandardCluster {
	c := NewCluster(prefix + "-master")
	c.RegisterSlave(prefix + "-slave01")
	c.RegisterSlave(prefix + "-slave02")
	c.RegisterSlave(prefix + "-slave03")
	return c
}

func TestShardClusterMaster(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	var node *Node

	s = &ShardCluster{slotsize: 1}
	node = s.Master()
	assert.Nil(node, "Master() should be always nil on ShardCluster")
}

func TestShardClusterMasters(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	var err error

	s = &ShardCluster{slotsize: 1000}
	assert.Len(s.Masters(), 0)

	c := testCreateCluster("shard01")
	err = s.RegisterShard(0, 500, c)
	assert.Nil(err)
	assert.Len(s.Masters(), 1)

	c = testCreateCluster("shard02")
	err = s.RegisterShard(501, 999, c)
	assert.Nil(err)
	assert.Len(s.Masters(), 2)
}

func TestShardClusterSlave(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	var node *Node

	s = &ShardCluster{slotsize: 1}
	node = s.Slave()
	assert.Nil(node, "Slave() should be always nil on ShardCluster")
}

func TestShardClusterSelectByKey(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	var c *StandardCluster
	var err error

	s = &ShardCluster{slotsize: 2}
	err = s.RegisterShard(0, 0, testCreateCluster("shard01"))
	assert.Nil(err)
	err = s.RegisterShard(1, 1, testCreateCluster("shard02"))
	assert.Nil(err)

	c = s.SelectByKey(0)
	assert.Equal("shard01-master", c.Master().DB())
	c = s.SelectByKey(1)
	assert.Equal("shard02-master", c.Master().DB())
	c = s.SelectByKey(2)
	assert.Equal("shard01-master", c.Master().DB())
	c = s.SelectByKey(3)
	assert.Equal("shard02-master", c.Master().DB())
	c = s.SelectByKey(4)
	assert.Equal("shard01-master", c.Master().DB())
	c = s.SelectByKey(5)
	assert.Equal("shard02-master", c.Master().DB())
}

func TestShardClusterRegisterShard(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	var err error

	s = &ShardCluster{slotsize: 10}
	err = s.RegisterShard(-1, 0, testCreateCluster("min-error"))
	assert.NotNil(err, "Slotsize cannot be under 0")
	assert.Len(s.List, 0)

	err = s.RegisterShard(0, 10, testCreateCluster("max-error"))
	assert.NotNil(err, "Slotsize cannot be greater equal than slotsize")
	assert.Len(s.List, 0)

	err = s.RegisterShard(5, 6, testCreateCluster("shard01"))
	assert.Nil(err)

	err = s.RegisterShard(6, 9, testCreateCluster("min-error"))
	assert.NotNil(err, "Slot min is already registered")

	err = s.RegisterShard(0, 5, testCreateCluster("max-error"))
	assert.NotNil(err, "Slot max is already registered")
}

func TestShardClusterCheckOverlapped(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	var err error

	s = &ShardCluster{slotsize: 10}
	err = s.RegisterShard(5, 6, testCreateCluster("shard01"))
	assert.Nil(err)

	err = s.checkOverlapped(6, 9)
	assert.NotNil(err, "Slot min is already registered")

	err = s.checkOverlapped(0, 5)
	assert.NotNil(err, "Slot max is already registered")

	err = s.checkOverlapped(0, 4)
	assert.Nil(err)
	err = s.checkOverlapped(7, 9)
	assert.Nil(err)
}

func TestShardSetInRange(t *testing.T) {
	assert := assert.New(t)

	var ss *ShardSet
	ss = &ShardSet{
		min: 10,
		max: 20,
	}

	assert.False(ss.InRange(9))
	assert.True(ss.InRange(10))
	assert.True(ss.InRange(11))
	assert.True(ss.InRange(19))
	assert.True(ss.InRange(20))
	assert.False(ss.InRange(21))
}

func TestShardSetCheckSlotSize(t *testing.T) {
	assert := assert.New(t)

	var ss *ShardSet
	var err error
	ss = &ShardSet{
		min: 10,
		max: 20,
	}

	err = ss.checkSlotSize(19)
	assert.NotNil(err, "max must be greater than slotsize")
	err = ss.checkSlotSize(20)
	assert.NotNil(err, "max must be greater than slotsize")
	err = ss.checkSlotSize(21)
	assert.Nil(err)

	ss = &ShardSet{
		min: -2,
		max: 20,
	}
	err = ss.checkSlotSize(21)
	assert.NotNil(err, "min must be greater equal than 0")
}

func TestShardSetIsMaxInSlotSize(t *testing.T) {
	assert := assert.New(t)

	var ss *ShardSet
	ss = &ShardSet{
		min: 10,
		max: 20,
	}

	assert.False(ss.isMaxInSlotSize(19))
	assert.False(ss.isMaxInSlotSize(20))
	assert.True(ss.isMaxInSlotSize(21))
}

func TestShardSetIsMinAboveZero(t *testing.T) {
	assert := assert.New(t)

	assert.False(ShardSet{min: -1}.isMinAboveZero())
	assert.True(ShardSet{min: 0}.isMinAboveZero())
	assert.True(ShardSet{min: 1}.isMinAboveZero())
}
