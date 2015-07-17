package xorm

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/evalphobia/wizard"
)

var w *wizard.Wizard
var engine1, engine2, engine3, engine4 *xorm.Engine

type User struct {
	ID   int64  `xorm:"pk not null" shard_key:"true"`
	Name string `xorm:"varchar(255) not null"`
}

func ExampleRegisterStandardDatabases() {
	engine1, _ = xorm.NewEngine("mysql", "root:@/example_user?charset=utf8")
	engine2, _ = xorm.NewEngine("mysql", "root:@/example_user?charset=utf8")
	engine3, _ = xorm.NewEngine("mysql", "root:@/example_foobar?charset=utf8")
	engine4, _ = xorm.NewEngine("mysql", "root:@/example_other?charset=utf8")

	w = wizard.NewWizard()
	stndardCluster := w.CreateCluster(User{}, engine1) // engine is master database used for table of User{}
	stndardCluster.RegisterSlave(engine2)              // add slave

	_ = w.CreateCluster("foobar", engine3) // engine3 is master database used for table of foobar

	stndardCluster = wizard.NewCluster(engine4)
	w.SetDefault(stndardCluster) // engine4 is master database used for all the other tables
}

func ExampleRegisterShardedDatabase() {
	engine1, _ = xorm.NewEngine("mysql", "root:@/example_user_a?charset=utf8")
	engine2, _ = xorm.NewEngine("mysql", "root:@/example_user_a?charset=utf8")
	engine3, _ = xorm.NewEngine("mysql", "root:@/example_user_b?charset=utf8")
	engine4, _ = xorm.NewEngine("mysql", "root:@/example_user_b?charset=utf8")

	w = wizard.NewWizard()
	shardClusters := w.CreateShardCluster(&User{}, 997) // create shard clusters for User{} with slotsize 997
	standardClusterA := wizard.NewCluster(engine1)
	standardClusterA.RegisterSlave(engine2)
	shardClusters.RegisterShard(0, 500, standardClusterA)

	standardClusterB := wizard.NewCluster(engine3)
	standardClusterB.RegisterSlave(engine4)
	shardClusters.RegisterShard(501, 996, standardClusterB)
}

func ExampleGet() {
	orm := New(w)

	user := &User{ID: 99}
	has, err := orm.Get(user, func(s Session) (bool, error) {
		return s.Get(user)
	})
	if err != nil {
		fmt.Printf("error occured, %s", err.Error())
		return
	}

	if !has {
		fmt.Printf("cannot find the user. id:%d", user.ID)
		return
	}
	fmt.Printf("user found. id:%d, name:%s", user.ID, user.Name)
}

func ExampleInsert() {
	orm := New(w)

	user := &User{ID: 99, Name: "Adam Smith"}
	total, err := orm.Insert(user, func(s Session) (int64, error) {
		return s.Insert(user)
	})
	if err != nil {
		fmt.Printf("error occured, %s", err.Error())
		return
	}
	if total < 1 {
		fmt.Printf("insert failed. id:%d", user.ID)
		return
	}
}

func ExampleTransaction() {
	var err error
	orm := New(w)

	user1 := &User{ID: 1, Name: "Adam Smith"}
	user2 := &User{ID: 2, Name: "Benjamin Franklin"}

	s1, _ := orm.Transaction(user1)
	s2, _ := orm.Transaction(user2)

	_, err = s1.Insert(user1)
	if err != nil {
		orm.RollbackAll()
		return
	}
	_, err = s2.Insert(user2)
	if err != nil {
		orm.RollbackAll()
		return
	}
	orm.CommitAll()
}

func ExampleTransactionAuto() {
	var err error
	orm := New(w)

	user1 := &User{ID: 1, Name: "Adam Smith"}
	user2 := &User{ID: 2, Name: "Benjamin Franklin"}

	orm.SetAutoTransaction(true)
	_, err = orm.Insert(user1, func(s Session) (int64, error) {
		return s.Insert(user1)
	})
	if err != nil {
		orm.RollbackAll()
		return
	}
	_, err = orm.Insert(user2, func(s Session) (int64, error) {
		return s.Insert(user2)
	})
	if err != nil {
		orm.RollbackAll()
		return
	}
	orm.CommitAll()
}
