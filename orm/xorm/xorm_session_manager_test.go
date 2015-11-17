package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testID = "my unique identifier"

func TestNewMasterSession(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{ID: 2}
	s, err = orm.NewMasterSession(row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{ID: 500}
	s, err = orm.NewMasterSession(row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)
}

func TestReadOnly(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	sl := orm.XormSessionManager.getOrCreateSessionList(testID)

	assert.False(sl.readOnly)
	orm.ReadOnly(testID, true)
	assert.True(sl.readOnly)
	orm.ReadOnly(testID, false)
	assert.False(sl.readOnly)
}

func TestIsReadOnly(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	assert.False(orm.IsReadOnly(testID))
	orm.ReadOnly(testID, true)
	assert.True(orm.IsReadOnly(testID))
	orm.ReadOnly(testID, false)
	assert.False(orm.IsReadOnly(testID))
}

func TestSetAutoTransaction(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	sl := orm.XormSessionManager.getOrCreateSessionList(testID)

	assert.False(sl.autoTx)
	orm.SetAutoTransaction(testID, true)
	assert.True(sl.autoTx)
	orm.SetAutoTransaction(testID, false)
	assert.False(sl.autoTx)
}

func TestSetIsAutoTransaction(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	assert.False(orm.IsAutoTransaction(testID))
	orm.SetAutoTransaction(testID, true)
	assert.True(orm.IsAutoTransaction(testID))
	orm.SetAutoTransaction(testID, false)
	assert.False(orm.IsAutoTransaction(testID))
}

func TestUseMasterSession(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{ID: 2}
	s, err = orm.UseMasterSession(testID, row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{ID: 500}
	s, err = orm.UseMasterSession(testID, row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// auto tx
	orm.SetAutoTransaction(testID, true)
	s, err = orm.UseMasterSession(testID, testUser{ID: 1})
	s.Insert(testUser{ID: 4})
	count, _ := s.Count(testUser{})
	assert.EqualValues(4, count)

	s.Rollback()
	s.Init()
	count, _ = s.Count(testUser{})
	assert.EqualValues(3, count)
}

func TestUseMasterSessionByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row *testUser
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{}
	s, err = orm.UseMasterSessionByKey(testID, row, 1)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 2
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{}
	s, err = orm.UseMasterSessionByKey(testID, row, 900)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 500
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// auto tx
	orm.SetAutoTransaction(testID, true)
	s, err = orm.UseMasterSessionByKey(testID, testUser{}, 1)
	s.Insert(testUser{ID: 4})
	count, _ := s.Count(testUser{})
	assert.EqualValues(4, count)

	s.Rollback()
	s.Init()
	count, _ = s.Count(testUser{})
	assert.EqualValues(3, count)
}

func TestUseSlaveSession(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{ID: 2}
	s, err = orm.UseSlaveSession(testID, row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{ID: 500}
	s, err = orm.UseSlaveSession(testID, row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)
}

func TestUseSlaveSessionByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row *testUser
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{}
	s, err = orm.UseSlaveSessionByKey(testID, row, 1)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 2
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{}
	s, err = orm.UseSlaveSessionByKey(testID, row, 900)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 500
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)
}
