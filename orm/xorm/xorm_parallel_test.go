package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindParallel(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error
	var list []testUser

	err = orm.FindParallel(&list, testUser{}, "id > ? ", 1)
	assert.Nil(err)
	assert.Len(list, 5)
	assert.Contains(list, testUser{ID: 2, Name: "Benjamin"})
	assert.Contains(list, testUser{ID: 3, Name: "Charles"})
	assert.Contains(list, testUser{ID: 500, Name: "Alice"})
	assert.Contains(list, testUser{ID: 501, Name: "Betty"})
	assert.Contains(list, testUser{ID: 502, Name: "Christina"})
}

func TestFindParallelByCondition(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error
	var list []testUser

	// order by asc
	cond := NewFindCondition(testUser{})
	cond.And("id > ?", 1)
	cond.SetLimit(1)
	cond.OrderByAsc("id")

	err = orm.FindParallelByCondition(&list, cond)
	assert.Nil(err)
	assert.Len(list, 2)
	assert.Contains(list, testUser{ID: 2, Name: "Benjamin"})
	assert.Contains(list, testUser{ID: 500, Name: "Alice"})

	// order by desc limit 1
	list = []testUser{}
	cond = NewFindCondition(testUser{})
	cond.And("id > ?", 1)
	cond.SetLimit(1)
	cond.OrderByDesc("id")

	err = orm.FindParallelByCondition(&list, cond)
	assert.Nil(err)
	assert.Len(list, 2)
	assert.Contains(list, testUser{ID: 3, Name: "Charles"})
	assert.Contains(list, testUser{ID: 502, Name: "Christina"})

	// order by desc limit 1 offset 1
	list = []testUser{}
	cond = NewFindCondition(testUser{})
	cond.And("id > ?", 1)
	cond.SetLimit(1)
	cond.OrderByDesc("id")
	cond.SetOffset(1)
	err = orm.FindParallelByCondition(&list, cond)
	assert.Nil(err)
	assert.Len(list, 2)
	assert.Contains(list, testUser{ID: 2, Name: "Benjamin"})
	assert.Contains(list, testUser{ID: 501, Name: "Betty"})
}

func TestCountParallelByCondition(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	var err error
	var testObj testUser

	// order by asc
	cond := NewFindCondition(testUser{})
	cond.And("id > ?", 1)

	counts, err := orm.CountParallelByCondition(&testObj, cond)
	assert.Nil(err)
	assert.Len(counts, 2)
	assert.Equal(counts[0], int64(3))
	assert.Equal(counts[1], int64(2))
}
