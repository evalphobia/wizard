package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeValue(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("wizard.ShardCluster", NormalizeValue(ShardCluster{}), "Struct should return the type name")
	assert.Equal("wizard.ShardCluster", NormalizeValue(&ShardCluster{}), "Struct pointer should return the type name")

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

func TestGetShardKey(t *testing.T) {
	assert := assert.New(t)

	type myStruct1 struct {
		UserID    int64
		CountryID int64
		CityID    int64 `shard_key:"false"`
	}

	type myStruct2 struct {
		UserID    int64
		CountryID int64 `shard_key:"true"`
		CityID    int64
	}

	type myStruct3 struct {
		UserID    int64 `shard_key:"true"`
		CountryID int64 `shard_key:"true"`
		CityID    int64 `shard_key:"true"`
	}

	type personStruct struct {
		Name string
		City string `shard_key:"true"`
		Tel  string
	}

	m1 := myStruct1{UserID: 1, CountryID: 2, CityID: 3}
	m2 := myStruct2{UserID: 1, CountryID: 2, CityID: 3}
	m3 := myStruct3{UserID: 1, CountryID: 2, CityID: 3}

	assert.Equal(int64(0), getShardKey(m1), "getShardKey() must return 0 when tag `shard_key:true` is missing")
	assert.Equal(m2.CountryID, getShardKey(m2))
	assert.Equal(m3.UserID, getShardKey(m3), "getShardKey() must return 1st field value when multiple tag `shard_key:true` exists")

	adam := personStruct{Name: "Adam Smith", City: "Oxford", Tel: "+81 0120-000-000"}
	assert.Equal(getInt64("Oxford"), getShardKey(adam))
}
