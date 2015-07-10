package wizard

// Node is struct for single database instance
type Node struct {
	db interface{} // db connection
}

// NewNode returns initialized Node
func NewNode(db interface{}) *Node {
	return &Node{db: db}
}

// DB is used for returning database connection
func (n *Node) DB() interface{} {
	return n.db
}
