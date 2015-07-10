package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAssertInt64(a *assert.Assertions, v interface{}) {
	a.Equal(int64(99), getInt64(v))
}

func TestGetInt64(t *testing.T) {
	a := assert.New(t)

	var vInt = 99
	testAssertInt64(a, vInt)

	var vInt8 int8 = 99
	testAssertInt64(a, vInt8)

	var vInt16 int16 = 99
	testAssertInt64(a, vInt16)

	var vInt32 int32 = 99
	testAssertInt64(a, vInt32)

	var vInt64 int64 = 99
	testAssertInt64(a, vInt64)

	var vUInt uint = 99
	testAssertInt64(a, vUInt)

	var vUInt8 uint8 = 99
	testAssertInt64(a, vUInt8)

	var vUInt16 uint16 = 99
	testAssertInt64(a, vUInt16)

	var vUInt32 uint32 = 99
	testAssertInt64(a, vUInt32)

	var vUInt64 uint64 = 99
	testAssertInt64(a, vUInt64)

	var vFloat32 float32 = 99
	testAssertInt64(a, vFloat32)

	var vFloat64 float64 = 99
	testAssertInt64(a, vFloat64)

	var vStr = "foobar"
	a.Equal(int64(3297785893580976128), getInt64(vStr))

	type myStruct struct{}
	a.Equal(int64(2612580365084131328), getInt64(myStruct{}))
}

func TestHashToInt64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(3297785893580976128), hashToInt64("foobar"))

	type myStruct struct{}
	assert.Equal(int64(2612580365084131328), hashToInt64(myStruct{}))
}
