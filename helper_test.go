package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseMaster(t *testing.T) {
	assert := assert.New(t)

	var wiz *Wizard
	var c *StandardCluster

	wiz = NewWizard()
	c = wiz.CreateCluster("country_table", "db-master")
	c.RegisterSlave("db-slave")

	assert.Equal("db-master", wiz.UseMaster("country_table"))
	assert.Nil(wiz.UseMaster("city_table"), "Non registered name")

	var s *ShardCluster
	s = wiz.CreateShardCluster("user_table", 997)
	s.RegisterShard(0, 499, NewCluster("shard01-master"))
	s.RegisterShard(500, 996, NewCluster("shard02-master"))
	assert.Equal("shard01-master", wiz.UseMaster("user_table"))
}

func TestUseMasters(t *testing.T) {
	assert := assert.New(t)

	var wiz *Wizard
	wiz = NewWizard()
	wiz.CreateCluster("country_table", "db-master")

	assert.Contains(wiz.UseMasters("country_table"), "db-master")
	assert.Empty(wiz.UseMasters("city_table"), "Non registered name")

	var s *ShardCluster
	s = wiz.CreateShardCluster("user_table", 997)
	s.RegisterShard(0, 499, NewCluster("shard01-master"))
	s.RegisterShard(500, 996, NewCluster("shard02-master"))

	assert.Contains(wiz.UseMasters("user_table"), "shard01-master")
	assert.Contains(wiz.UseMasters("user_table"), "shard02-master")
	assert.Len(wiz.UseMasters("user_table"), 2)
}

func TestUseSlave(t *testing.T) {
	assert := assert.New(t)

	var wiz *Wizard
	var c *StandardCluster

	wiz = NewWizard()
	c = wiz.CreateCluster("country_table", "db-master")
	assert.Equal("db-master", wiz.UseSlave("country_table"), "Slave() return master when no slaves exists")

	c.RegisterSlave("db-slave")
	assert.Equal("db-slave", wiz.UseSlave("country_table"))

	assert.Nil(wiz.UseSlave("city_table"), "Non registered name")
}

func TestUseMasterBySlot(t *testing.T) {
	assert := assert.New(t)

	var wiz *Wizard
	var c *StandardCluster

	wiz = NewWizard()
	c = wiz.CreateCluster("country_table", "db-master")
	c.RegisterSlave("db-slave")

	assert.Equal("db-master", wiz.UseMasterBySlot("country_table", 1))
	assert.Nil(wiz.UseMasterBySlot("city_table", 1), "Non registered name")

	var s *ShardCluster
	s = wiz.CreateShardCluster("user_table", 997)
	s.RegisterShard(0, 499, NewCluster("shard01-master"))
	s.RegisterShard(500, 996, NewCluster("shard02-master"))
	assert.Equal("shard01-master", wiz.UseMasterBySlot("user_table", 499))
	assert.Equal("shard02-master", wiz.UseMasterBySlot("user_table", 500))
	assert.Equal("shard02-master", wiz.UseMasterBySlot("user_table", 996))
	assert.Equal("shard01-master", wiz.UseMasterBySlot("user_table", 997))
}

// TODO: add test for multiple slaves
func TestUseSlaveBySlot(t *testing.T) {
	assert := assert.New(t)

	var wiz *Wizard
	var c *StandardCluster

	wiz = NewWizard()
	c = wiz.CreateCluster("country_table", "db-master")
	c.RegisterSlave("db-slave")

	assert.Equal("db-master", wiz.UseMasterBySlot("country_table", 1))
	assert.Nil(wiz.UseMasterBySlot("city_table", 1), "Non registered name")

	var s *ShardCluster
	s = wiz.CreateShardCluster("user_table", 997)
	c1 := NewCluster("shard01-master")
	c1.RegisterSlave("shard01-slave")
	c2 := NewCluster("shard02-master")
	c2.RegisterSlave("shard02-slave")
	s.RegisterShard(0, 499, c1)
	s.RegisterShard(500, 996, c2)
	assert.Equal("shard01-slave", wiz.UseSlaveBySlot("user_table", 499))
	assert.Equal("shard02-slave", wiz.UseSlaveBySlot("user_table", 500))
	assert.Equal("shard02-slave", wiz.UseSlaveBySlot("user_table", 996))
	assert.Equal("shard01-slave", wiz.UseSlaveBySlot("user_table", 997))
}
