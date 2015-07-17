package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	// auto tx
	orm.SetAutoTransaction(true)
	s, err = orm.NewMasterSession(testUser{ID: 1})
	s.Insert(testUser{ID: 4})
	count, _ := s.Count(testUser{})
	assert.EqualValues(4, count)

	s.Rollback()
	s.Init()
	count, _ = s.Count(testUser{})
	assert.EqualValues(3, count)
}

func TestNewMasterSessionByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row *testUser
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{}
	s, err = orm.NewMasterSessionByKey(row, 1)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 2
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{}
	s, err = orm.NewMasterSessionByKey(row, 900)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 500
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// auto tx
	orm.SetAutoTransaction(true)
	s, err = orm.NewMasterSessionByKey(testUser{}, 1)
	s.Insert(testUser{ID: 4})
	count, _ := s.Count(testUser{})
	assert.EqualValues(4, count)

	s.Rollback()
	s.Init()
	count, _ = s.Count(testUser{})
	assert.EqualValues(3, count)
}

func TestNewSlaveSession(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row interface{}
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{ID: 2}
	s, err = orm.NewSlaveSession(row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{ID: 500}
	s, err = orm.NewSlaveSession(row)
	assert.Nil(err)
	assert.NotNil(s)

	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)
}

func TestNewSlaveSessionByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var row *testUser
	var s Session
	var has bool
	var err error

	// A
	row = &testUser{}
	s, err = orm.NewSlaveSessionByKey(row, 1)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 2
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)

	// B
	row = &testUser{}
	s, err = orm.NewSlaveSessionByKey(row, 900)
	assert.Nil(err)
	assert.NotNil(s)

	row.ID = 500
	has, err = s.Get(row)
	assert.Nil(err)
	assert.True(has)
}
