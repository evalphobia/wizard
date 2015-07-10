package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	assert := assert.New(t)

	var n *Node
	n = NewNode("db")
	assert.Equal("db", n.db, "db should be saved on NewNode()")

	type TestDB struct{}
	n = NewNode(TestDB{})
	assert.Equal(TestDB{}, n.db, "db should be saved on NewNode()")
}

func TestDB(t *testing.T) {
	assert := assert.New(t)

	var n *Node
	n = NewNode("db")
	assert.Equal("db", n.DB(), "DB() should equal to Node.db")

	type TestDB struct{}
	n = NewNode(TestDB{})
	assert.Equal(TestDB{}, n.DB(), "db should equal to Node.db")
}
