package xorm

import (
	"database/sql"

	"github.com/go-xorm/xorm"
)

type Session interface{
	Close()
	Sql(string, ...interface{}) *xorm.Session
	Where(string, ...interface{}) *xorm.Session
	And(string, ...interface{}) *xorm.Session
	Or(string, ...interface{}) *xorm.Session
	
	Id(interface{}) *xorm.Session
	Table(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Select(string) *xorm.Session
	Cols(...string) *xorm.Session
	MustCols(...string) *xorm.Session
	AllCols() *xorm.Session
	Distinct(...string) *xorm.Session
	Omit(...string) *xorm.Session
	Nullable(...string) *xorm.Session
	NoAutoTime() *xorm.Session

	Limit(int, ...int) *xorm.Session
	OrderBy(string) *xorm.Session
	Desc(...string) *xorm.Session
	Asc(...string) *xorm.Session
	Join(string, interface{}, string) *xorm.Session
	GroupBy(string) *xorm.Session
	Having(string) *xorm.Session
	
	Begin() error
	Rollback() error
	Commit( ) error
	Exec(string, ...interface{}) (sql.Result, error) 
	
	CreateTable(interface{}) error
	CreateIndexes(interface{}) error
	CreateUniques(interface{}) error
	DropIndexes(interface{}) error
	DropTable(interface{}) error
	
	Rows(interface{}) (*xorm.Rows, error)
	Iterate(interface{}, xorm.IterFunc) error
	Get(interface{}) (bool, error)
	Count(interface{}) (int64, error)
	Find(interface{}, ...interface{}) error
	
	IsTableExist(interface{}) (bool, error)
	IsTableEmpty(interface{}) (bool, error)

	Query(string, ...interface{}) ([]map[string][]byte, error)
	Insert(...interface{}) (int64, error)
	InsertMulti(interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Update(interface{}, ...interface{}) (int64, error)
	Delete(interface{}) (int64, error)
	
	Sync2(...interface{}) error
}

