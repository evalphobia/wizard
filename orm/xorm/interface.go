package xorm

import (
	"database/sql"
	"io"
	"time"

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
	Init()
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
	ForUpdate() *xorm.Session
	Omit(...string) *xorm.Session
	Nullable(...string) *xorm.Session
	NoAutoTime() *xorm.Session
	Incr(string, ...interface{}) *xorm.Session
	Decr(string, ...interface{}) *xorm.Session
	SetExpr(string, string) *xorm.Session

	Limit(int, ...int) *xorm.Session
	OrderBy(string) *xorm.Session
	Desc(...string) *xorm.Session
	Asc(...string) *xorm.Session
	Join(string, interface{}, string) *xorm.Session
	GroupBy(string) *xorm.Session
	Having(string) *xorm.Session

	Begin() error
	Rollback() error
	Commit() error
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

	LastSQL() (string, []interface{})
	Sync2(...interface{}) error
}

// Engine is interface for xorm.Engine
type Engine interface {
	SetDisableGlobalCache(bool)
	DriverName() string
	DataSourceName() string
	SetMapper(core.IMapper)
	SetTableMapper(core.IMapper)
	SetColumnMapper(core.IMapper)
	SupportInsertMany() bool
	QuoteStr() string
	Quote(string) string
	AutoIncrStr() string
	SetMaxOpenConns(int)
	SetMaxConns(int)
	SetMaxIdleConns(int)
	NoCache() *xorm.Session
	NoCascade() *xorm.Session
	SetLogger(core.ILogger)

	NewSession() *xorm.Session
	Close() error
	Ping() error

	Sql(string, ...interface{}) *xorm.Session
	NoAutoTime() *xorm.Session
	DumpAllToFile(string) error
	DumpAll(io.Writer) error

	Cascade(...bool) *xorm.Session
	Where(string, ...interface{}) *xorm.Session
	Id(interface{}) *xorm.Session

	Before(func(interface{})) *xorm.Session
	After(func(interface{})) *xorm.Session
	Charset(string) *xorm.Session
	StoreEngine(string) *xorm.Session

	Distinct(...string) *xorm.Session
	Select(string) *xorm.Session
	Cols(...string) *xorm.Session
	AllCols() *xorm.Session
	MustCols(...string) *xorm.Session
	UseBool(...string) *xorm.Session
	Omit(...string) *xorm.Session
	Nullable(...string) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Incr(string, ...interface{}) *xorm.Session
	Decr(string, ...interface{}) *xorm.Session

	Table(interface{}) *xorm.Session
	Limit(int, ...int) *xorm.Session
	Desc(...string) *xorm.Session
	Asc(...string) *xorm.Session
	OrderBy(string) *xorm.Session
	Join(string, interface{}, string) *xorm.Session
	GroupBy(string) *xorm.Session
	Having(string) *xorm.Session

	GobRegister(interface{}) *xorm.Engine

	IsTableEmpty(interface{}) (bool, error)
	IsTableExist(interface{}) (bool, error)
	CreateIndexes(interface{}) error
	CreateUniques(interface{}) error
	ClearCacheBean(interface{}, string) error
	ClearCache(...interface{}) error
	Sync(...interface{}) error
	Sync2(...interface{}) error
	CreateTables(...interface{}) error
	DropTables(...interface{}) error

	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) ([]map[string][]byte, error)
	Insert(...interface{}) (int64, error)
	Update(interface{}, ...interface{}) (int64, error)
	Delete(interface{}) (int64, error)
	Get(interface{}) (bool, error)
	Find(interface{}, ...interface{}) error
	Iterate(interface{}, xorm.IterFunc) error
	Rows(interface{}) (*xorm.Rows, error)
	Count(interface{}) (int64, error)

	ImportFile(string) ([]sql.Result, error)
	Import(io.Reader) ([]sql.Result, error)

	TZTime(time.Time) time.Time
	NowTime(string) interface{}
	NowTime2(string) (interface{}, time.Time)
	FormatTime(string, time.Time) interface{}
	Unscoped() *xorm.Session
}
