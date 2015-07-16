package wizard

import (
	"reflect"
	"strings"
)

// used for shard key in the tag name of struct
const TagName = "shard_key"

// NormalizeValue returns value
// if struct is passed, returns name of the struct
// if pointer is passed, returns non-pointer value
func NormalizeValue(p interface{}) interface{} {
	v := toValue(p)
	if v.Kind() == reflect.Struct {
		return v.Type().String()
	}
	return v.Interface()
}

func getShardKey(p interface{}) int64 {
	v := toValue(p)
	if v.Kind() != reflect.Struct {
		return 0
	}
	return getShardKeyFromStruct(p, TagName)
}

// toValue converts any value to reflect.Value
func toValue(p interface{}) reflect.Value {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// toType converts any value to reflect.Type
func toType(p interface{}) reflect.Type {
	t := reflect.ValueOf(p).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func getShardKeyFromStruct(p interface{}, tagName string) int64 {
	t := toType(p)
	values := toValue(p)
	for i, max := 0, t.NumField(); i < max; i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		tag := parseTag(f, tagName)
		if tag != "true" {
			continue
		}
		v := values.Field(i)
		return getInt64(v.Interface())
	}
	return 0
}

// parseTag returns the first tag value of the struct field
func parseTag(f reflect.StructField, tag string) string {
	res := strings.Split(f.Tag.Get(tag), ",")
	return res[0]
}
