package xorm

import (
	"github.com/evalphobia/wizard"
)

func GetBySession(s Session, fn func(Session) (bool, error)) (bool, error) {
	return fn(s)
}

func FindBySession(s Session, fn func(Session) error) error {
	return fn(s)
}

func CountBySession(s Session, fn func(Session) (int, error)) (int, error) {
	return fn(s)
}

func InsertBySession(s Session, fn func(Session) (int64, error)) (int64, error) {
	return fn(s)
}

func UpdateBySession(s Session, fn func(Session) (bool, error)) (bool, error) {
	return fn(s)
}

func NormalizeValue(p interface{}) interface{} {
	return wizard.NormalizeValue(p)
}