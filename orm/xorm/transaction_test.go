package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	initTestDB()
}

func countUserSlave(orm *Xorm) int64 {
	count, _ := orm.Count(&testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Count(&testUser{})
	})
	return count
}

func countUserMaster(orm *Xorm) int64 {
	count, _ := orm.CountUsingMaster(&testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Count(&testUser{})
	})
	return count
}

func TestBeginSession(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	s, err := orm.BeginSession(testUser{ID: 1})
	assert.Nil(err)
	assert.NotNil(s)

	count, _ := s.Count(testUser{})
	assert.EqualValues(3, count)

	s.Insert(&testUser{ID: 4})
	count, _ = s.Count(testUser{})
	assert.EqualValues(4, count)

	err = s.Rollback()
	assert.Nil(err)

	s, _ = orm.BeginSession(testUser{ID: 1})
	count, _ = s.Count(testUser{})
	assert.EqualValues(3, count)
}

func TestBegin(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error
	err = orm.Begin(&testUser{ID: 1})
	assert.Nil(err)

	assert.EqualValues(3, countUserSlave(orm))

	orm.Insert(&testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Insert(&testUser{ID: 4, Name: "Daniel"})
	})
	assert.EqualValues(3, countUserSlave(orm))  // slave used
	assert.EqualValues(4, countUserMaster(orm)) // master used
	initTestDB()
}

func TestCommit(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error
	err = orm.Begin(&testUser{ID: 1})
	assert.Nil(err)

	assert.EqualValues(3, countUserSlave(orm))

	orm.Insert(&testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Insert(&testUser{ID: 4, Name: "Daniel"})
	})
	assert.EqualValues(3, countUserSlave(orm))

	orm.Commit(&testUser{ID: 1})
	assert.EqualValues(4, countUserSlave(orm))

	initTestDB()
}

func TestRollback(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error
	err = orm.Begin(&testUser{ID: 1})
	assert.Nil(err)

	assert.EqualValues(3, countUserMaster(orm))

	orm.Insert(&testUser{ID: 1}, func(s Session) (int64, error) {
		return s.Insert(&testUser{ID: 4, Name: "Daniel"})
	})
	assert.EqualValues(4, countUserMaster(orm))

	err = orm.Rollback(&testUser{ID: 1})
	assert.Nil(err)
	assert.EqualValues(3, countUserMaster(orm))
	initTestDB()
}
