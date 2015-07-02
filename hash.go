package wizard

import (
	"fmt"
	"hash/crc64"
	"io"
)

var hashTable = crc64.MakeTable(crc64.ISO)

func getInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float32:
		return int64(t)
	case float64:
		return int64(t)
	}
	return hashToInt64(v)
}

func hashToInt64(v interface{}) int64 {
	str := fmt.Sprint(v)
	h := crc64.New(hashTable)
	io.WriteString(h, str)
	return int64(h.Sum64())
}
