[![Build Status](https://drone.io/github.com/evalphobia/wizard/status.png)](https://drone.io/github.com/evalphobia/wizard/latest)

[![Coverage Status](https://coveralls.io/repos/evalphobia/wizard/badge.svg?branch=master&service=github)](https://coveralls.io/github/evalphobia/wizard?branch=master)

# Wizard

[![Join the chat at https://gitter.im/evalphobia/wizard](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/evalphobia/wizard?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Wizard is database/sql management library for multi instance and sharding in golang.  
Inspired by [MixedGauge](https://github.com/taiki45/mixed_gauge)

## Supported orm list

- [xorm](https://github.com/go-xorm/xorm)

## Quick Usage

### Register database clusters

```go
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/evalphobia/wizard"
)

type Blog struct {
	ArticleID string `xorm:"article_id pk VARCHAR(255) not null"`
	Content   string `xorm:"content text"`
}

type User struct {
	ID   int64  `xorm:"id pk BIGINT(20) not null" shard_key:"true"`
	Name string `xorm:"name VARCHAR(100) not null"`
}

func main() {
	wiz = wizard.NewWizard()

	/**
		register normal cluster
	*/

	// create engines
	blogMaster, _ = xorm.NewEngine("mysql", "root:@tcp(db-master:3306)/blog?charset=utf8")
	blogSlave01, _ = xorm.NewEngine("mysql", "root:@tcp(db-slave01:3306)/blog?charset=utf8")
	blogSlave02, _ = xorm.NewEngine("mysql", "root:@tcp(db-slave01:3306)/blog?charset=utf8")

	// create cluster with master nodel; CreateCluster(name, master-instance)
	blogCluster := wiz.CreateCluster(Blog{}, blogMaster)
	blogCluster.RegisterSlave(blogSlave01) // add slaves
	blogCluster.RegisterSlave(blogSlave02)


	/**
		register shard clusters
	*/

	// shard one
	user01Master, _ = xorm.NewEngine("mysql", "root:@/tcp(shard01-master:3306)/users?charset=utf8")
	user01Slave01, _ = xorm.NewEngine("mysql", "root:@/tcp(shard01-slave01:3306)/users?charset=utf8")
	user01Slave02, _ = xorm.NewEngine("mysql", "root:@/tcp(shard01-slave02:3306)/users?charset=utf8")

	// shard two
	user02Master, _ = xorm.NewEngine("mysql", "root:@/tcp(shard02-master:3306)/users?charset=utf8")
	user02Slave01, _ = xorm.NewEngine("mysql", "root:@/tcp(shard02-slave01:3306)/users?charset=utf8")
	user02Slave02, _ = xorm.NewEngine("mysql", "root:@/tcp(shard02-slave02:3306)/users?charset=utf8")

	// create shard clusters; CreateShardCluster(name, slot-size)
	shardClusters := wiz.CreateShardCluster(User{}, 1023)

	// create single shard set #1
	shardCluster01 := wizard.NewCluster(user01Master)
	shardCluster01.RegisterSlave(user01Slave01)
	shardCluster01.RegisterSlave(user01Slave02)

	// create single shard set #2
	shardCluster02 := wizard.NewCluster(user02Master)
	shardCluster02.RegisterSlave(user02Slave01)
	shardCluster02.RegisterSlave(user02Slave02)

	// register shards with slot; RegisterShard(min, max, cluster)
	shardClusters.RegisterShard(0, 500, shardCluster01)
	shardClusters.RegisterShard(501, 1022, shardCluster02)
}
```

### Query on database clusters

```go
import (
	"fmt"

	"github.com/evalphobia/wizard/orm/xorm"
)

func main() {
	orm := xorm.New(wiz)

	blog := &Blog{ArticleID: "hello-world"}
	has, err := orm.Get(blog, func(s xorm.Session) (bool, error) {
		return s.Get(blog)
	})
	// => SELECT * FROM blog WHERE article_id = "hello-world"; -- execute on blog SLAVE


	fmt.Println(has)
	fmt.Println(blog.Content)

	user := &User{
		ID:   1600, // 1600 % 1023 = 577; => shard02
		Name: "Adam Smith",
	}

	err = orm.Begin(user) // => BEGIN; -- execute on user02-MASTER
	if err != nil {
		panic("Error on transaction beginning")
	}

	total, err := orm.Insert(user, func(s xorm.Session) (int64, error) {
		return s.Insert(user)
	})
	// => INSERT INTO users VALUES(1600, "Adam Smith"); -- execute on user02-MASTER

	fmt.Println(total)

	newUser := &User{ID: 1600}
	has, err = orm.GetUsingMaster(newUser, func(s xorm.Session) (bool, error) {
		return s.Get(newUser)
	})
	// => SELECT * FROM users WHERE id = 1600; -- execute on user02-MASTER


	err = orm.Commit(user)
	// => COMMIT; -- execute on user02-MASTER
	if err != nil {
		panic("Error on transaction ending")
	}
}
```

### Notes

- Clusters is selected by name, which can be any value like `string`, `struct`, `pointer`.
    - the pointer value automatically converts to the non-pointer value.
- Struct field tag: `shard_key:"true"` is used as a shard-key
    - shard_key is divided by slot size and the mod value is used for shard mapping
    - string shard_key convert to int64 with CRC64 and divided by slot-size
