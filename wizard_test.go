package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWizard(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	assert.NotNil(wiz)
	assert.Empty(wiz.clusters)
}

func TestSetDefault(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	assert.Nil(wiz.defaultCluster)

	c := NewCluster("db")
	wiz.SetDefault(c)
	assert.Equal(c, wiz.defaultCluster, "It can set StandardCluster")

	s := &ShardCluster{}
	wiz.SetDefault(s)
	assert.Equal(s, wiz.defaultCluster, "It can set ShardCluster")
}

func TestHasDefault(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	assert.False(wiz.HasDefault())

	c := NewCluster("db")
	wiz.SetDefault(c)
	assert.True(wiz.HasDefault())
}

func TestGetCluster(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	assert.Nil(wiz.getCluster("table name"))

	c := NewCluster("db")
	wiz.clusters["table name"] = c
	assert.Equal(c, wiz.getCluster("table name"))
}

func TestSetCluster(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	assert.Nil(wiz.clusters["table name"])

	c := NewCluster("db")
	wiz.setCluster(c, "table name")
	assert.Equal(c, wiz.clusters["table name"])
}

func TestCreateCluster(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	c := wiz.CreateCluster("table name", "db-master")
	assert.NotNil(c)
	assert.NotNil(c.master)
	assert.Empty(c.slaves)
	assert.Equal("db-master", c.master.db)
}

func TestCreateShardCluster(t *testing.T) {
	assert := assert.New(t)

	var s *ShardCluster
	wiz := NewWizard()

	var slotsize int64 = 99
	s = wiz.CreateShardCluster("table name", slotsize)
	assert.NotNil(s)
	assert.Empty(s.List)
	assert.Equal(int64(99), s.slotsize)
	assert.Equal(s, wiz.clusters["table name"])

	var slotsizeZero int64
	s = wiz.CreateShardCluster("table name2", slotsizeZero)
	assert.NotNil(s)
	assert.Empty(s.List)
	assert.Equal(int64(1), s.slotsize)
	assert.Equal(s, wiz.clusters["table name2"])

	var slotsizeMinus int64 = -99
	s = wiz.CreateShardCluster("table name", slotsizeMinus)
	assert.NotNil(s)
	assert.Empty(s.List)
	assert.Equal(int64(1), s.slotsize)
	assert.Equal(s, wiz.clusters["table name"])
}

func TestSelect(t *testing.T) {
	assert := assert.New(t)
	wiz := NewWizard()

	// shard test
	type myStruct struct {
		ID int64 `shard_key:"true"`
	}
	s := wiz.CreateShardCluster(myStruct{}, 100)
	shardSet1 := NewCluster("shard01-master")
	shardSet2 := NewCluster("shard02-master")
	s.RegisterShard(0, 49, shardSet1)
	s.RegisterShard(50, 99, shardSet2)

	assert.Equal(shardSet1, wiz.Select(&myStruct{ID: 1}))
	assert.Equal(shardSet1, wiz.Select(&myStruct{ID: 49}))
	assert.Equal(shardSet2, wiz.Select(&myStruct{ID: 50}))
	assert.Equal(shardSet2, wiz.Select(&myStruct{ID: 99}))
	assert.Equal(shardSet1, wiz.Select(&myStruct{ID: 100}))
	assert.Equal(shardSet1, wiz.Select(&myStruct{ID: 149}))
	assert.Equal(shardSet2, wiz.Select(&myStruct{ID: 150}))
	assert.Equal(shardSet2, wiz.Select(&myStruct{ID: 199}))
	assert.Equal(shardSet1, wiz.Select(&myStruct{ID: 200}))

	// standard test
	c1 := wiz.CreateCluster("standard table", "db-master")
	c2 := wiz.Select("standard table")
	assert.Equal(c1, c2)

	c3 := wiz.CreateCluster(myStruct{}, "db-master")
	c4 := wiz.Select(&myStruct{ID: 100})
	assert.Equal(c3, c4)

	// error
	_ = wiz.CreateShardCluster("shard table", 100)
	nilShard := wiz.Select("shard table")
	assert.Nil(nilShard, "Select() returns nil for shardcluster when obj does not contain shardkey")

	nilTable := wiz.Select("not registered")
	assert.Nil(nilTable, "Select() returns nil when table name does not registered")
}

func TestSelectByKey(t *testing.T) {
	assert := assert.New(t)

	wiz := NewWizard()
	c1 := wiz.CreateCluster("standard table", "db-master")
	c2 := wiz.SelectByKey("standard table", 1)
	c3 := wiz.SelectByKey("standard table", 99)
	assert.Equal(c1, c2)
	assert.Equal(c1, c3)

	// object test
	type myStruct struct {
		ID int64 `shard_key:"true"`
	}
	s1 := wiz.CreateShardCluster(myStruct{}, 100)
	shardSet1 := NewCluster("shard01-master")
	shardSet2 := NewCluster("shard02-master")
	s1.RegisterShard(0, 49, shardSet1)
	s1.RegisterShard(50, 99, shardSet2)

	assert.Equal(shardSet1, wiz.SelectByKey(&myStruct{ID: 99}, 1))
	assert.Equal(shardSet1, wiz.SelectByKey(&myStruct{ID: 99}, 49))
	assert.Equal(shardSet2, wiz.SelectByKey(&myStruct{ID: 99}, 50))
	assert.Equal(shardSet2, wiz.SelectByKey(&myStruct{ID: 99}, 99))
	assert.Equal(shardSet1, wiz.SelectByKey(&myStruct{ID: 99}, 100))
	assert.Equal(shardSet1, wiz.SelectByKey(&myStruct{ID: 99}, 149))
	assert.Equal(shardSet2, wiz.SelectByKey(&myStruct{ID: 99}, 150))
	assert.Equal(shardSet2, wiz.SelectByKey(&myStruct{ID: 99}, 199))
	assert.Equal(shardSet1, wiz.SelectByKey(&myStruct{ID: 99}, 200))

	// non object test
	s2 := wiz.CreateShardCluster("shard table", 100)
	s2.RegisterShard(0, 49, NewCluster("x01-master"))
	s2.RegisterShard(50, 99, NewCluster("x02-master"))
	c4 := s2.SelectByKey(5000)
	c5 := wiz.SelectByKey("shard table", 5000)
	assert.Equal(c4, c5, "Select() returns nil for shardcluster when obj does not contain shardkey")

	// error
	nilTable := wiz.SelectByKey("not registered", 99)
	assert.Nil(nilTable, "Select() returns nil when table name does not registered")
}
