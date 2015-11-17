package xorm

import (
	"os"
	"testing"
	"time"

	"github.com/evalphobia/wizard"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	dbUser01Master, dbUser01Slave01, dbUser01Slave02 Engine // user A
	dbUser02Master, dbUser02Slave01, dbUser02Slave02 Engine // user B
	dbFoobarMaster, dbFoobarSlave01, dbFoobarSlave02 Engine
	dbOther                                          Engine
	wiz, emptyWiz                                    *wizard.Wizard
)

type testUser struct {
	ID   int64  `xorm:"id pk not null" shard_key:"true"`
	Name string `xorm:"varchar(255) not null"`
}

func (u testUser) TableName() string {
	return "test_user"
}

type testFoobar struct {
	ID   int64  `xorm:"id pk not null"`
	Name string `xorm:"varchar(255) not null"`
}

func (f testFoobar) TableName() string {
	return "test_foobar"
}

type testCompany struct {
	ID   int64  `xorm:"id pk not null"`
	Name string `xorm:"varchar(255) not null"`
}

func (c testCompany) TableName() string {
	return "test_company"
}

func init() {
	initTestDB()
	emptyWiz = wizard.NewWizard()
}

func initTestDB() {
	testInitializeEngines()
	testWaitForIO()
	testInitializeSchema()
	testWaitForIO()
	testInitializeData()
	testWaitForIO()
}

func testWaitForIO() {
	time.Sleep(80 * time.Millisecond)
}

func testInitializeEngines() {
	f1 := "xorm_test_user01.db"
	f2 := "xorm_test_user02.db"
	f3 := "xorm_test_foobar.db"
	f4 := "xorm_test_other.db"
	os.Remove(f1)
	os.Remove(f2)
	os.Remove(f3)
	os.Remove(f4)

	dbUser01Master, _ = xorm.NewEngine("sqlite3", f1)
	dbUser01Slave01, _ = xorm.NewEngine("sqlite3", f1)
	dbUser01Slave02, _ = xorm.NewEngine("sqlite3", f1)
	dbUser02Master, _ = xorm.NewEngine("sqlite3", f2)
	dbUser02Slave01, _ = xorm.NewEngine("sqlite3", f2)
	dbUser02Slave02, _ = xorm.NewEngine("sqlite3", f2)
	dbFoobarMaster, _ = xorm.NewEngine("sqlite3", f3)
	dbFoobarSlave01, _ = xorm.NewEngine("sqlite3", f3)
	dbFoobarSlave02, _ = xorm.NewEngine("sqlite3", f3)
	dbOther, _ = xorm.NewEngine("sqlite3", f4)
}

func testInitializeSchema() {
	dbUser01Master.Sync(&testUser{})
	dbUser02Master.Sync(&testUser{})
	dbFoobarMaster.Sync(&testFoobar{})
	dbOther.Sync(&testCompany{})
}

func testInitializeData() {
	dbUser01Master.Delete(testUser{})
	dbUser02Master.Delete(testUser{})
	dbFoobarMaster.Delete(testFoobar{})
	dbOther.Delete(testCompany{})

	dbUser01Master.Insert(testUser{ID: 1, Name: "Adam"})
	dbUser01Master.Insert(testUser{ID: 2, Name: "Benjamin"})
	dbUser01Master.Insert(testUser{ID: 3, Name: "Charles"})
	dbUser02Master.Insert(testUser{ID: 500, Name: "Alice"})
	dbUser02Master.Insert(testUser{ID: 501, Name: "Betty"})
	dbUser02Master.Insert(testUser{ID: 502, Name: "Christina"})
	dbFoobarMaster.Insert(testFoobar{ID: 1, Name: "foobar#1"})
	dbFoobarMaster.Insert(testFoobar{ID: 2, Name: "foobar#2"})
	dbFoobarMaster.Insert(testFoobar{ID: 3, Name: "foobar#3"})
	dbOther.Insert(testCompany{ID: 1, Name: "Apple"})
	dbOther.Insert(testCompany{ID: 2, Name: "BOX"})
	dbOther.Insert(testCompany{ID: 3, Name: "Criteo"})
}

func testCreateWizard() *wizard.Wizard {
	wiz := wizard.NewWizard()

	userShards := wiz.CreateShardCluster(testUser{}, 997)
	shard01 := wizard.NewCluster(dbUser01Master)
	shard01.RegisterSlave(dbUser01Slave01)
	shard01.RegisterSlave(dbUser01Slave02)
	userShards.RegisterShard(0, 499, shard01) // user A

	shard02 := wizard.NewCluster(dbUser02Master)
	shard02.RegisterSlave(dbUser02Slave01)
	shard02.RegisterSlave(dbUser02Slave02)
	userShards.RegisterShard(500, 996, shard02) // user B

	foobarCluster := wiz.CreateCluster(testFoobar{}, dbFoobarMaster)
	foobarCluster.RegisterSlave(dbFoobarSlave01)
	foobarCluster.RegisterSlave(dbFoobarSlave02)

	otherCluster := wizard.NewCluster(dbOther)
	wiz.SetDefault(otherCluster)
	return wiz
}

func countUserMaster(orm *Xorm) int64 {
	count, _ := orm.CountUsingMaster(testID, &testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Count(&testUser{})
	})
	return count
}

func countUserMasterB(orm *Xorm) int64 {
	count, _ := orm.CountUsingMaster(testID, &testUser{ID: 500}, func(s Session) (int64, error) {
		return s.Count(&testUser{})
	})
	return count
}

func countUserSlave(orm *Xorm) int64 {
	count, _ := orm.Count(&testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Count(&testUser{})
	})
	return count
}

func countUserSlaveB(orm *Xorm) int64 {
	count, _ := orm.Count(&testUser{ID: 500}, func(s Session) (int64, error) {
		return s.Count(&testUser{})
	})
	return count
}

func countUserBySession(s Session) int64 {
	count, _ := s.Count(testUser{})
	return count
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	wiz := wizard.NewWizard()

	orm := New(wiz)
	assert.Equal(wiz, orm.Wiz)
}
