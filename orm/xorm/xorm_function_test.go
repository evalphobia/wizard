package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		s.And("id = 1")
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

	testWaitForIO()

	// user A
	assert.EqualValues(3, countFn(testUser{ID: 1}, &testUser{}))
	row = &testUser{ID: 1000, Name: "Daniel"}
	affected, err = orm.Insert(testID, row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testUser{ID: 1}, &testUser{}))

	row = &testUser{ID: 1000}
	orm.Get(row, getFn)
	assert.Equal("Daniel", row.(*testUser).Name)

	testWaitForIO()

	// user B
	assert.EqualValues(3, countFn(testUser{ID: 500}, &testUser{}))
	row = &testUser{ID: 1500, Name: "Dorothy"}
	affected, err = orm.Insert(testID, row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testUser{ID: 500}, &testUser{}))

	row = &testUser{ID: 1500}
	orm.Get(row, getFn)
	assert.Equal("Dorothy", row.(*testUser).Name)

	testWaitForIO()

	// foobar
	assert.EqualValues(3, countFn(testFoobar{}, &testFoobar{}))
	row = &testFoobar{ID: 4, Name: "foobar#4"}
	affected, err = orm.Insert(testID, row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testFoobar{}, &testFoobar{}))

	row = &testFoobar{ID: 4}
	orm.Get(row, getFn)
	assert.Equal("foobar#4", row.(*testFoobar).Name)

	testWaitForIO()

	// other
	assert.EqualValues(3, countFn(testCompany{}, &testCompany{}))
	row = &testCompany{ID: 4, Name: "Delta Air Lines"}
	affected, err = orm.Insert(testID, row, fn)
	assert.Nil(err)
	assert.Equal(success, affected)
	assert.EqualValues(4, countFn(testCompany{}, &testCompany{}))

	row = &testCompany{ID: 4}
	orm.Get(row, getFn)
	assert.Equal("Delta Air Lines", row.(*testCompany).Name)

	testWaitForIO()

	// multiple rows
	assert.EqualValues(4, countFn(testCompany{}, &testCompany{}))
	rows := []*testCompany{
		{ID: 5, Name: "eureka"},
		{ID: 6, Name: "Facebook"},
		{ID: 7, Name: "Google"},
	}
	affected, err = orm.Insert(testID, testCompany{}, func(s Session) (int64, error) {
		return s.Insert(&rows)
	})
	assert.Nil(err)
	assert.Equal(int64(3), affected)
	assert.EqualValues(7, countFn(testCompany{}, &testCompany{}))

	testWaitForIO()

	// readonly
	orm.ReadOnly(testID, true)
	affected, err = orm.Insert(testID, row, fn)
	assert.Nil(err)
	assert.EqualValues(0, affected)

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
	affected, err = orm.Update(testID, user, func(s Session) (int64, error) {
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
	affected, err = orm.Update(testID, user, func(s Session) (int64, error) {
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
	affected, err = orm.Update(testID, foobar, func(s Session) (int64, error) {
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
	affected, err = orm.Update(testID, company, func(s Session) (int64, error) {
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
	affected, err = orm.Update(testID, foobar, func(s Session) (int64, error) {
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

	// readonly
	orm.ReadOnly(testID, true)
	affected, err = orm.Update(testID, foobar, func(s Session) (int64, error) {
		return s.Update(foobar)
	})
	assert.Nil(err)
	assert.EqualValues(0, affected)

	initTestDB()
}

func TestGetUsingMaster(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var has bool
	var err error
	fn := func(s Session) (bool, error) {
		return s.Get(row)
	}

	orm.SetAutoTransaction(testID, true)

	row = &testUser{ID: 4}
	s, _ := orm.Transaction(testID, row)
	s.Insert(row)

	// slave
	has, err = orm.Get(row, fn)
	assert.Nil(err)
	assert.False(has)

	// master
	has, err = orm.GetUsingMaster(testID, row, fn)
	assert.Nil(err)
	assert.True(has)

	s.Rollback()
}

func TestFindUsingMaster(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var foobars []testFoobar
	var err error
	fn := func(s Session) error {
		return s.Find(&foobars)
	}

	orm.SetAutoTransaction(testID, true)

	row := &testFoobar{ID: 4, Name: "foobar#4@FindUsingMaster"}
	s, _ := orm.Transaction(testID, row)
	s.Insert(row)

	// slave
	err = orm.Find(testFoobar{}, fn)
	assert.Nil(err)
	assert.Len(foobars, 3)

	// master
	foobars = foobars[:0]
	err = orm.FindUsingMaster(testID, testFoobar{}, fn)
	assert.Nil(err)
	assert.Len(foobars, 4)

	s.Rollback()
}

func TestCountUsingMaster(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var foobar = testFoobar{}
	var count int64
	var err error
	fn := func(s Session) (int64, error) {
		return s.Count(&foobar)
	}

	orm.SetAutoTransaction(testID, true)

	row := &testFoobar{ID: 4, Name: "foobar#4@CountUsingMaster"}
	s, _ := orm.Transaction(testID, row)
	s.Insert(row)

	count, err = orm.Count(foobar, fn)
	assert.Nil(err)
	assert.EqualValues(3, count)

	count, err = orm.CountUsingMaster(testID, foobar, fn)
	assert.Nil(err)
	assert.EqualValues(4, count)

	s.Rollback()
}

func TestFunctionNilDB(t *testing.T) {
	assert := assert.New(t)
	orm := New(emptyWiz)

	// var nilFoobars []*testFoobar
	var has bool
	var count, affected int64
	var err error

	fnGet := func(s Session) (bool, error) { return true, nil }
	fnFind := func(s Session) error { return nil }
	fnCount := func(s Session) (int64, error) { return 99, nil }

	// Get
	has, err = orm.Get(testFoobar{}, fnGet)
	assert.NotNil(err)
	assert.False(has)

	// Find
	err = orm.Find(testFoobar{}, fnFind)
	assert.NotNil(err)

	// Count
	count, err = orm.Count(testFoobar{}, fnCount)
	assert.NotNil(err)
	assert.EqualValues(0, count)

	// GetUsingMaster
	has, err = orm.GetUsingMaster(testID, testFoobar{}, fnGet)
	assert.NotNil(err)
	assert.False(has)

	// FindUsingMaster
	err = orm.FindUsingMaster(testID, testFoobar{}, fnFind)
	assert.NotNil(err)

	// CountUsingMaster
	count, err = orm.CountUsingMaster(testID, testFoobar{}, fnCount)
	assert.NotNil(err)
	assert.EqualValues(0, count)

	// Insert
	affected, err = orm.Insert(testID, testFoobar{}, fnCount)
	assert.NotNil(err)
	assert.EqualValues(0, affected)

	// nil db
	affected, err = orm.Update(testID, testFoobar{}, fnCount)
	assert.NotNil(err)
	assert.EqualValues(0, affected)
}
