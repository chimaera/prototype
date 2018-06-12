package db

import (
	"errors"
	"reflect"
	"sync"
)

type Node struct {
	sync.RWMutex

	Type      NodeType
	Value     interface{}
	Relations []*Node
}

func NewNode(t NodeType, val interface{}) *Node {
	return &Node{
		Type:      t,
		Value:     val,
		Relations: make([]*Node, 0),
	}
}

func (n *Node) Equals(v *Node) bool {
	n.RLock()
	defer n.RUnlock()

	if n.Type != v.Type {
		return false
	}

	return reflect.DeepEqual(n.Value, v.Value)
}

func (n *Node) Add(t NodeType, v interface{}) {
	n.Connect(NewNode(t, v))
}

func (n *Node) Connect(v *Node) error {
	if !n.IsConnected(v) {
		n.Lock()
		defer n.Unlock()
		n.Relations = append(n.Relations, v)
		return nil
	} else {
		return errors.New("node connect: node already connected")
	}
}

func (n *Node) Find(v *Node) *Node {
	n.RLock()
	defer n.RUnlock()

	if n.Equals(v) {
		return n
	}

	for _, r := range n.Relations {
		if r.Equals(v) {
			return r
		} else if found := r.Find(v); found != nil {
			return found
		}
	}

	return nil
}

func (n *Node) IsConnected(v *Node) bool {
	return n.Find(v) != nil
}

type PrinterFN func(format string, v ...interface{})

func doPad(fn PrinterFN, pad int) {
	for i := 0; i < pad; i++ {
		fn("  ")
	}
}

func (n *Node) Print(fn PrinterFN, pad int) {
	n.RLock()
	defer n.RUnlock()

	doPad(fn, pad)
	fn("%s[%v]\n", n.Type.String(), n.Value)

	for _, r := range n.Relations {
		r.Print(fn, pad+1)
	}
}
