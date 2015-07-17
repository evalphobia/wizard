package xorm

import (
	"github.com/evalphobia/wizard"
)

// GetBySession performs wrapped function for xorm.Sesion.Get()
func GetBySession(s Session, fn func(Session) (bool, error)) (bool, error) {
	return fn(s)
}

// FindBySession performs wrapped function for xorm.Sesion.Find()
func FindBySession(s Session, fn func(Session) error) error {
	return fn(s)
}

// CountBySession performs wrapped function for xorm.Sesion.Count()
func CountBySession(s Session, fn func(Session) (int, error)) (int, error) {
	return fn(s)
}

// InsertBySession performs wrapped function for xorm.Sesion.Insert()
func InsertBySession(s Session, fn func(Session) (int64, error)) (int64, error) {
	return fn(s)
}

// UpdateBySession performs wrapped function for xorm.Sesion.Update()
func UpdateBySession(s Session, fn func(Session) (bool, error)) (bool, error) {
	return fn(s)
}

// NormalizeValue returns non-pointer value
func NormalizeValue(p interface{}) interface{} {
	return wizard.NormalizeValue(p)
}
