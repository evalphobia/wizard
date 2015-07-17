package xorm

import (
	"testing"

	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
)

func TestGetBySession(t *testing.T) {
	assert := assert.New(t)
	sess := &xorm.Session{}

	fn := func(s Session) (bool, error) {
		assert.Equal(sess, s)
		return true, nil
	}
	b, err := GetBySession(sess, fn)
	assert.True(b)
	assert.Nil(err)
}

func TestFindBySession(t *testing.T) {
	assert := assert.New(t)
	sess := &xorm.Session{}

	fn := func(s Session) error {
		assert.Equal(sess, s)
		return nil
	}
	err := FindBySession(sess, fn)
	assert.Nil(err)
}

func TestCountBySession(t *testing.T) {
	assert := assert.New(t)
	sess := &xorm.Session{}

	fn := func(s Session) (int, error) {
		assert.Equal(sess, s)
		return 99, nil
	}
	count, err := CountBySession(sess, fn)
	assert.Equal(99, count)
	assert.Nil(err)
}

func TestInsertBySession(t *testing.T) {
	assert := assert.New(t)
	sess := &xorm.Session{}

	fn := func(s Session) (int64, error) {
		assert.Equal(sess, s)
		return 99, nil
	}
	affected, err := InsertBySession(sess, fn)
	assert.Equal(int64(99), affected)
	assert.Nil(err)
}

func TestUpdateBySession(t *testing.T) {
	assert := assert.New(t)
	sess := &xorm.Session{}

	fn := func(s Session) (bool, error) {
		assert.Equal(sess, s)
		return true, nil
	}
	b, err := UpdateBySession(sess, fn)
	assert.True(b)
	assert.Nil(err)
}

func TestNormalizeValue(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("xorm.Xorm", NormalizeValue(Xorm{}), "Struct should return the type name")
	assert.Equal("xorm.Xorm", NormalizeValue(&Xorm{}), "Struct pointer should return the type name")

	valueString := "foobar"
	assert.Equal(valueString, NormalizeValue(valueString))
	assert.Equal(valueString, NormalizeValue(&valueString))

	valueInt := 99
	assert.Equal(valueInt, NormalizeValue(valueInt))
	assert.Equal(valueInt, NormalizeValue(&valueInt))

	valueSlice := []string{"a", "b", "c"}
	assert.Equal(valueSlice, NormalizeValue(valueSlice))
	assert.Equal(valueSlice, NormalizeValue(&valueSlice))

	valueMap := map[interface{}]interface{}{"key": "value", 100: 403}
	assert.Equal(valueMap, NormalizeValue(valueMap))
	assert.Equal(valueMap, NormalizeValue(&valueMap))
}
