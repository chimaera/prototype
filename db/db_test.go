package db

import "testing"

func TestDatabase(t *testing.T) {
	db := Database{}

	if db.root != nil {
		t.Error("expected root of the database to be nil")
	}
}

func TestDatabaseCurrent(t *testing.T) {
	if Current != nil {
		t.Error("expected to start off with the current database as nil")
	}
}

func TestNewDatabase(t *testing.T) {
	rootNode := &Node{}

	db := New(rootNode)

	if db.root != rootNode {
		t.Error("unable to create new DB with a given root node")
	}
}

func TestDatabaseSearch(t *testing.T) {
	rootNode := &Node{Type: NodeTypeInfo, Value: "Example"}

	db := New(rootNode)

	result := db.Search(NodeTypeInfo, "Example")

	if result.Value != "Example" {
		t.Error("unable to retrieve node from search")
	}
}

func TestDatabaseRoot(t *testing.T) {
	rootNode := &Node{Type: NodeTypeInfo, Value: "Example"}

	db := New(rootNode)

	if db.Root() != rootNode {
		t.Error("unable to retrieve root node")
	}
}
