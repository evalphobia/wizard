package wizard

type Node struct {
	db interface{}
}

func NewNode(db interface{}) *Node {
	return &Node{db: db}
}

func (n *Node) DB() interface{} {
	return n.db
}
