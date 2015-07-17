package xorm

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMaster(t *testing.T) {
	assert := assert.New(t)

	wiz := testCreateWizard()
	orm := New(wiz)

	assert.Equal(dbFoobarMaster, orm.Master(testFoobar{}))
	assert.Equal(dbOther, orm.Master("xxx"))

	emptyOrm := New(emptyWiz)
	assert.Equal(nil, emptyOrm.Master("empty"))
}

func TestMasterByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	assert.Equal(dbUser01Master, orm.MasterByKey(testUser{}, 1))
	assert.Equal(dbUser01Master, orm.MasterByKey(testUser{}, 499))
	assert.Equal(dbUser02Master, orm.MasterByKey(testUser{}, 500))
	assert.Equal(dbUser02Master, orm.MasterByKey(testUser{}, 501))
	assert.Equal(dbUser02Master, orm.MasterByKey(testUser{}, 996))
	assert.Equal(dbUser01Master, orm.MasterByKey(testUser{}, 997))

	emptyOrm := New(emptyWiz)
	assert.Equal(nil, emptyOrm.MasterByKey("empty", 1))
}

func TestMasters(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	shardMasters := []Engine{dbUser01Master, dbUser02Master}
	assert.Equal(shardMasters, orm.Masters(testUser{}))
	assert.Equal([]Engine{dbFoobarMaster}, orm.Masters(testFoobar{}))
	assert.Equal([]Engine{dbOther}, orm.Masters("xxx"))

	emptyOrm := New(emptyWiz)
	assert.Empty(emptyOrm.Masters("empty"))
}

func TestSlave(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	assert.Contains([]Engine{dbFoobarSlave01, dbFoobarSlave02}, orm.Slave(testFoobar{}))
	assert.Equal(dbOther, orm.Slave("xxx"))

	emptyOrm := New(emptyWiz)
	assert.Equal(nil, emptyOrm.Slave("empty"))
}

func TestSlaveByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	db01 := []Engine{dbUser01Slave01, dbUser01Slave02}
	db02 := []Engine{dbUser02Slave01, dbUser02Slave02}
	assert.Contains(db01, orm.SlaveByKey(testUser{}, 1))
	assert.Contains(db01, orm.SlaveByKey(testUser{}, 499))
	assert.Contains(db02, orm.SlaveByKey(testUser{}, 500))
	assert.Contains(db02, orm.SlaveByKey(testUser{}, 501))
	assert.Contains(db02, orm.SlaveByKey(testUser{}, 996))
	assert.Contains(db01, orm.SlaveByKey(testUser{}, 997))

	emptyOrm := New(emptyWiz)
	assert.Equal(nil, emptyOrm.SlaveByKey("empty", 1))
}

func TestSlaves(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	slaves := orm.Slaves(testFoobar{})
	assert.Contains([]Engine{dbFoobarSlave01, dbFoobarSlave02}, slaves[0])

	slaves = orm.Slaves("xxx")
	assert.Equal([]Engine{dbOther}, slaves)

	slaves = orm.Slaves(testUser{})
	assert.Len(slaves, 2)
	assert.Contains([]Engine{dbUser01Slave01, dbUser01Slave02}, slaves[0])
	assert.Contains([]Engine{dbUser02Slave01, dbUser02Slave02}, slaves[1])

	emptyOrm := New(emptyWiz)
	assert.Empty(emptyOrm.Slaves("empty"))
}
