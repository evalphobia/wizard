package xorm

import (
	"database/sql"
	"io"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// ORM is wrapper interface for wizard.Xorm
type ORM interface {
	ReadOnly(Identifier, bool)
	IsReadOnly(Identifier) bool
	SetAutoTransaction(Identifier, bool)
	IsAutoTransaction(Identifier) bool

	Master(interface{}) Engine
	MasterByKey(interface{}, interface{}) Engine
	Masters(interface{}) []Engine
	Slave(interface{}) Engine
	SlaveByKey(interface{}, interface{}) Engine
	Slaves(interface{}) []Engine

	Get(interface{}, func(Session) (bool, error)) (bool, error)
	Find(interface{}, func(Session) error) error
	Count(interface{}, func(Session) (int64, error)) (int64, error)
	Insert(Identifier, interface{}, func(Session) (int64, error)) (int64, error)
	Update(Identifier, interface{}, func(Session) (int64, error)) (int64, error)
	FindParallel(interface{}, interface{}, string, ...interface{}) error
	FindParallelByCondition(interface{}, FindCondition) error
	CountParallelByCondition(interface{}, FindCondition) ([]int64, error)
	UpdateParallelByCondition(interface{}, UpdateCondition) (int64, error)
	GetUsingMaster(Identifier, interface{}, func(Session) (bool, error)) (bool, error)
	FindUsingMaster(Identifier, interface{}, func(Session) error) error
	CountUsingMaster(Identifier, interface{}, func(Session) (int64, error)) (int64, error)

	NewMasterSession(interface{}) (Session, error)

	UseMasterSession(Identifier, interface{}) (Session, error)
	UseMasterSessionByKey(Identifier, interface{}, interface{}) (Session, error)
	UseSlaveSession(Identifier, interface{}) (Session, error)
	UseSlaveSessionByKey(Identifier, interface{}, interface{}) (Session, error)
	UseAllMasterSessions(Identifier, interface{}) ([]Session, error)

	ForceNewTransaction(interface{}) (Session, error)
	Transaction(Identifier, interface{}) (Session, error)
	TransactionByKey(Identifier, interface{}, interface{}) (Session, error)
	AutoTransaction(Identifier, interface{}, Session) error
	CommitAll(Identifier) error
	RollbackAll(Identifier) error
	CloseAll(Identifier)
}

// Session is interface for xorm.Session
type Session interface {
	xorm.Interface

	And(interface{}, ...interface{}) *xorm.Session
	Begin() error
	Close()
	Commit() error
	CreateTable(interface{}) error
	DropTable(interface{}) error
	ForUpdate() *xorm.Session
	Having(string) *xorm.Session
	Id(interface{}) *xorm.Session
	Init()
	InsertMulti(interface{}) (int64, error)
	LastSQL() (string, []interface{})
	NoAutoTime() *xorm.Session
	Nullable(...string) *xorm.Session
	Or(interface{}, ...interface{}) *xorm.Session
	Rollback() error
	Select(string) *xorm.Session
	Sql(string, ...interface{}) *xorm.Session
	Sync2(...interface{}) error
}

// Engine is interface for xorm.Engine
type Engine interface {
	xorm.EngineInterface

	After(func(interface{})) *xorm.Session
	AutoIncrStr() string
	Cascade(...bool) *xorm.Session
	ClearCacheBean(interface{}, string) error
	Close() error
	DataSourceName() string
	DriverName() string
	DumpAll(io.Writer, ...core.DbType) error
	GobRegister(interface{}) *xorm.Engine
	Having(string) *xorm.Session
	Id(interface{}) *xorm.Session
	Import(io.Reader) ([]sql.Result, error)
	ImportFile(string) ([]sql.Result, error)
	NoCache() *xorm.Session
	NoCascade() *xorm.Session
	Nullable(...string) *xorm.Session
	QuoteStr() string
	Select(string) *xorm.Session
	SetColumnMapper(core.IMapper)
	SetDisableGlobalCache(bool)
	SetTableMapper(core.IMapper)
	Sql(string, ...interface{}) *xorm.Session
	SupportInsertMany() bool
}

var _ Session = &xorm.Session{}
var _ Engine = &xorm.Engine{}
