package db

type Database struct {
	root *Node
}

var (
	Current = (*Database)(nil)
)

func New(root *Node) *Database {
	Current = &Database{
		root: root,
	}
	return Current
}

func (db *Database) Search(t NodeType, v interface{}) (node *Node) {
	return db.root.Find(NewNode(t, v))
}

func (db *Database) SearchOrRoot(t NodeType, v interface{}) (node *Node) {
	if node = db.root.Find(NewNode(t, v)); node == nil {
		node = db.root
	}
	return
}

func (db *Database) Root() *Node {
	return db.root
}
