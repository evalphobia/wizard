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
