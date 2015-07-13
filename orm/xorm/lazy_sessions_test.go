package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	initTestDB()
}

func TestNewLazySessionList(t *testing.T) {
	assert := assert.New(t)

	s := newLazySessionList()
	assert.NotNil(s)
	assert.Empty(s.list)
	assert.Empty(s.inTx)

}

func TestNewLazySessions(t *testing.T) {
	assert := assert.New(t)

	ls := newLazySessions()
	assert.NotNil(ls)
	assert.Empty(ls.sessions)
}

func TestGetOrNewSession(t *testing.T) {
	assert := assert.New(t)

	ls := newLazySessions()

	// create new one
	newSession, err := ls.getOrNewSession(dbFoobarMaster)
	assert.Nil(err)
	assert.NotNil(newSession)

	savedSession, ok := ls.sessions[dbFoobarMaster]
	assert.True(ok)
	assert.Equal(savedSession, newSession)

	// from cache
	newSession2, err := ls.getOrNewSession(dbFoobarMaster)
	assert.Nil(err)
	assert.Equal(newSession, newSession2)
}

func TestLazySessionCommitAll(t *testing.T) {
	assert := assert.New(t)

	ls := newLazySessions()
	newSession, _ := ls.getOrNewSession(dbFoobarMaster)
	newSession.Insert(&testFoobar{ID: 4, Name: "foobar#4"})

	// before commit in another tx
	newSession2 := dbFoobarMaster.NewSession()
	count, err := newSession2.Count(&testFoobar{})
	assert.Nil(err)
	assert.EqualValues(3, count)

	err = ls.CommitAll()
	assert.Nil(err)

	// after commit in another tx
	count, err = newSession2.Count(&testFoobar{})
	assert.Nil(err)
	assert.EqualValues(4, count)

	initTestDB()
}

func TestLazySessionRollbackAll(t *testing.T) {
	assert := assert.New(t)

	ls := newLazySessions()
	newSession, _ := ls.getOrNewSession(dbFoobarMaster)
	newSession.Insert(&testFoobar{ID: 4, Name: "foobar#4"})

	// before rollback
	newSession, _ = ls.getOrNewSession(dbFoobarMaster)
	count, err := newSession.Count(&testFoobar{})
	assert.Nil(err)
	assert.EqualValues(4, count)

	err = ls.RollbackAll()
	assert.Nil(err)

	// after rollback
	newSession, _ = ls.getOrNewSession(dbFoobarMaster)
	count, err = newSession.Count(&testFoobar{})
	assert.Nil(err)
	assert.EqualValues(3, count)

	initTestDB()
}
