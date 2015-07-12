package xorm

import (
	"os"
	"testing"

	"github.com/evalphobia/wizard"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	dbUser01Master, dbUser01Slave01, dbUser01Slave02 *xorm.Engine // user A
	dbUser02Master, dbUser02Slave01, dbUser02Slave02 *xorm.Engine // user B
	dbFoobarMaster, dbFoobarSlave01, dbFoobarSlave02 *xorm.Engine
	dbOther                                          *xorm.Engine
	wiz                                              *wizard.Wizard
)

func init() {
	initTestDB()
}

func initTestDB() {
	testInitializeEngines()
	testInitializeSchema()
	testInitializeData()
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

func TestNew(t *testing.T) {
	assert := assert.New(t)
	wiz := wizard.NewWizard()

	orm := New(wiz)
	assert.Equal(wiz, orm.c)
}

func TestUseMaster(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()

	orm := New(wiz)
	assert.Equal(dbFoobarMaster, orm.UseMaster(testFoobar{}))
	assert.Equal(dbOther, orm.UseMaster("xxx"))
}

func TestUseMasters(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	shardMasters := []*xorm.Engine{dbUser01Master, dbUser02Master}
	assert.Equal(shardMasters, orm.UseMasters(testUser{}))
	assert.Equal([]*xorm.Engine{dbFoobarMaster}, orm.UseMasters(testFoobar{}))
	assert.Equal([]*xorm.Engine{dbOther}, orm.UseMasters("xxx"))
}

func TestUseSlave(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	assert.Contains([]*xorm.Engine{dbFoobarSlave01, dbFoobarSlave02}, orm.UseSlave(testFoobar{}))
	assert.Equal(dbOther, orm.UseSlave("xxx"))
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var has bool
	var err error
	fn := func(s Session) (bool, error) {
		return s.Get(row)
	}

	row = &testUser{ID: 1}
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.True(has)
	assert.Equal(int64(1), row.(*testUser).ID)
	assert.Equal("Adam", row.(*testUser).Name)

	row = &testUser{ID: 501}
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.True(has)
	assert.Equal(int64(501), row.(*testUser).ID)
	assert.Equal("Betty", row.(*testUser).Name)

	row = &testFoobar{ID: 1}
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.True(has)
	assert.Equal(int64(1), row.(*testFoobar).ID)
	assert.Equal("foobar#1", row.(*testFoobar).Name)

	row = &testCompany{ID: 2}
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.True(has)
	assert.Equal(int64(2), row.(*testCompany).ID)
	assert.Equal("BOX", row.(*testCompany).Name)

	// not found
	row = &testUser{ID: 4}
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.False(has)
	assert.Equal(int64(4), row.(*testUser).ID)
	assert.Equal("", row.(*testUser).Name)

	// not found
	row = &testUser{ID: 504}
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.False(has)
	assert.Equal(int64(504), row.(*testUser).ID)
	assert.Equal("", row.(*testUser).Name)
}

func TestFind(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error

	// user A
	var usersA []*testUser
	err = orm.Find(testUser{}, func(s Session) error {
		s.Where("id > 1")
		return s.Find(&usersA)
	})
	assert.Nil(err)
	assert.Len(usersA, 2)

	// user B
	var usersB []*testUser
	err = orm.Find(testUser{ID: 500}, func(s Session) error {
		s.Where("id > 1")
		return s.Find(&usersB)
	})
	assert.Nil(err)
	assert.Len(usersB, 3)

	var foobars []*testFoobar
	err = orm.Find(&testFoobar{}, func(s Session) error {
		s.Id(1)
		return s.Find(&foobars)
	})
	assert.Nil(err)
	assert.Len(foobars, 1)

	var companies []*testCompany
	err = orm.Find(&testCompany{}, func(s Session) error {
		s.Where("id > 2")
		return s.Find(&companies)
	})
	assert.Nil(err)
	assert.Len(companies, 1)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var count int64
	var err error

	// user A
	count, err = orm.Count(&testUser{ID: 1}, func(s Session) (int64, error) {
		s.Where("id > 1")
		return s.Count(&testUser{})
	})
	assert.Nil(err)
	assert.EqualValues(2, count)

	// user B
	count, err = orm.Count(&testUser{ID: 501}, func(s Session) (int64, error) {
		s.Where("id > 1")
		return s.Count(&testUser{})
	})
	assert.Nil(err)
	assert.EqualValues(3, count)

	count, err = orm.Count(&testFoobar{}, func(s Session) (int64, error) {
		return s.Count(&testFoobar{ID: 1})
	})
	assert.Nil(err)
	assert.EqualValues(1, count)

	count, err = orm.Count(&testCompany{}, func(s Session) (int64, error) {
		s.Where("id > 2")
		return s.Count(&testCompany{})
	})
	assert.Nil(err)
	assert.EqualValues(1, count)
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var affected int64
	var err error
	var success int64 = 1
	fn := func(s Session) (int64, error) {
		return s.Insert(row)
	}
	getFn := func(s Session) (bool, error) {
		return s.Get(row)
	}
	countFn := func(table, obj interface{}) int64 {
		count, _ := orm.Count(table, func(s Session) (int64, error) {
			return s.Count(obj)
		})
		return count
	}

	// user A
	assert.EqualValues(3, countFn(testUser{ID: 1}, &testUser{}))
	row = &testUser{ID: 1000, Name: "Daniel"}
	affected, err = orm.Insert(row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testUser{ID: 1}, &testUser{}))

	row = &testUser{ID: 1000}
	orm.Get(row, getFn)
	assert.Equal("Daniel", row.(*testUser).Name)

	// user B
	assert.EqualValues(3, countFn(testUser{ID: 500}, &testUser{}))
	row = &testUser{ID: 1500, Name: "Dorothy"}
	affected, err = orm.Insert(row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testUser{ID: 500}, &testUser{}))

	row = &testUser{ID: 1500}
	orm.Get(row, getFn)
	assert.Equal("Dorothy", row.(*testUser).Name)

	// foobar
	assert.EqualValues(3, countFn(testFoobar{}, &testFoobar{}))
	row = &testFoobar{ID: 4, Name: "foobar#4"}
	affected, err = orm.Insert(row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testFoobar{}, &testFoobar{}))

	row = &testFoobar{ID: 4}
	orm.Get(row, getFn)
	assert.Equal("foobar#4", row.(*testFoobar).Name)

	// other
	assert.EqualValues(3, countFn(testCompany{}, &testCompany{}))
	row = &testCompany{ID: 4, Name: "Delta Air Lines"}
	affected, err = orm.Insert(row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testCompany{}, &testCompany{}))

	row = &testCompany{ID: 4}
	orm.Get(row, getFn)
	assert.Equal("Delta Air Lines", row.(*testCompany).Name)

	// multiple rows
	assert.EqualValues(4, countFn(testCompany{}, &testCompany{}))
	rows := []*testCompany{
		&testCompany{ID: 5, Name: "eureka"},
		&testCompany{ID: 6, Name: "Facebook"},
		&testCompany{ID: 7, Name: "Google"},
	}
	affected, err = orm.Insert(testCompany{}, func(s Session) (int64, error) {
		return s.Insert(&rows)
	})
	assert.Nil(err)
	assert.Equal(int64(3), affected)
	assert.EqualValues(7, countFn(testCompany{}, &testCompany{}))

	initTestDB()
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var affected int64
	var err error
	var success int64 = 1
	getFn := func(s Session) (bool, error) {
		return s.Get(row)
	}

	// user A
	var user *testUser
	user = &testUser{ID: 1, Name: "Akira"}
	affected, err = orm.Update(user, func(s Session) (int64, error) {
		s.Where("id = ?", user.ID)
		return s.Update(user)
	})
	assert.Nil(err)
	assert.Equal(success, affected)

	row = &testUser{ID: 1}
	orm.Get(row, getFn)
	assert.Equal("Akira", row.(*testUser).Name)

	// // user B
	user = &testUser{ID: 501, Name: "Aiko"}
	affected, err = orm.Update(user, func(s Session) (int64, error) {
		s.Where("id = ?", user.ID)
		return s.Update(user)
	})
	assert.Nil(err)
	assert.Equal(success, affected)

	row = &testUser{ID: 501}
	orm.Get(row, getFn)
	assert.Equal("Aiko", row.(*testUser).Name)

	// foobar
	var foobar *testFoobar
	foobar = &testFoobar{ID: 1, Name: "foobar#1b"}
	affected, err = orm.Update(foobar, func(s Session) (int64, error) {
		s.Where("id = ?", foobar.ID)
		return s.Update(foobar)
	})
	assert.Nil(err)
	assert.Equal(success, affected)

	row = &testFoobar{ID: 1}
	orm.Get(row, getFn)
	assert.Equal("foobar#1b", row.(*testFoobar).Name)

	// other
	var company *testCompany
	company = &testCompany{ID: 1, Name: "Alibaba"}
	affected, err = orm.Update(company, func(s Session) (int64, error) {
		s.Where("id = ?", company.ID)
		return s.Update(company)
	})
	assert.Nil(err)
	assert.Equal(success, affected)

	row = &testCompany{ID: 1}
	orm.Get(row, getFn)
	assert.Equal("Alibaba", row.(*testCompany).Name)

	// multiple rows
	foobar = &testFoobar{Name: "foobar#XXX"}
	affected, err = orm.Update(foobar, func(s Session) (int64, error) {
		return s.Update(foobar)
	})
	assert.Nil(err)
	assert.Equal(int64(3), affected)

	row = &testFoobar{ID: 1}
	orm.Get(row, getFn)
	assert.Equal("foobar#XXX", row.(*testFoobar).Name)
	row = &testFoobar{ID: 2}
	orm.Get(row, getFn)
	assert.Equal("foobar#XXX", row.(*testFoobar).Name)
	row = &testFoobar{ID: 3}
	orm.Get(row, getFn)
	assert.Equal("foobar#XXX", row.(*testFoobar).Name)

	initTestDB()
}

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
