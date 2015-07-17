package xorm

import (
	"os"
	"testing"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func TestInterface(t *testing.T) {
	wiz := testCreateWizard()

	var orm ORM
	orm = New(wiz)
	_ = orm

	var e Engine
	name := "test_if.db"
	e, _ = xorm.NewEngine("sqlite3", name)
	os.Remove(name)

	var s Session
	s = e.NewSession()
	_ = s
}
