package db

import "testing"
import "reflect"

func TestNode(t *testing.T) {
	// bring some nodes to the party
	n1 := &Node{Type: NodeTypeInfo, Value: "Node 1"}
	n2 := &Node{Type: NodeTypeInfo, Value: "Node 2"}
	n3 := &Node{Type: NodeTypeInfo, Value: "Node 3"}

	// connect with node1 like we linkedin bruh
	n1.Connect(n2)
	n1.Connect(n3)

	// testing table
	var units = []struct {
		got interface{} // what we go
		exp interface{} // what we expect
		dsc interface{} // small description of error
	}{
		{n1.Type, NodeTypeInfo, "have a type"},
		{n1.Value, "Node 1", "have a value"},
		{len(n1.Relations), 2, "have one relation"},
		{len(n2.Relations), 0, "have no relations"},
		{n1.Find(n1), n1, "find node (self)"},
		{n1.Find(n2), n2, "find relational node"},
		{n1.IsConnected(n2), true, "be a known connection"},
		{n1.Equals(n1), true, "equals itself"},
		{n1.Equals(n2), false, "not equals different node"},
	}

	for _, u := range units {
		if !reflect.DeepEqual(u.exp, u.got) {
			t.Fatalf("expected '%v', got '%v', node should %v", u.exp, u.got, u.dsc)
		}
	}

}
