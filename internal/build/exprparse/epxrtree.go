package exprparse

import (
	"fmt"
)

type ParseTree struct {
	cur    *ParseTree
	parent *ParseTree
	Nodes  []Node
}

type Node struct {
	Val       string
	InnerTree *ParseTree
}

func NewRefTree(parent *ParseTree) *ParseTree {
	newTree := ParseTree{
		Nodes:  []Node{},
		parent: parent,
	}
	newTree.cur = &newTree
	return &newTree
}

func (r *ParseTree) AddField(field string) error {
	r.AddNode(Node{Val: field})
	return nil
}

func (r *ParseTree) AddNode(n Node) {
	r.cur.Nodes = append(r.cur.Nodes, n)
}

func (r *ParseTree) CloseInner() error {
	if r.parent == nil {
		return fmt.Errorf("closing subexpression without opening")
	}
	r.parent = r.parent.parent
	r.cur = r.parent
	return nil
}

func (r *ParseTree) OpenInner() {
	InnerTree := NewRefTree(r)
	r.AddNode(Node{InnerTree: InnerTree})
	// set cur to inner until bracket is closed
	r.parent = r.cur
	r.cur = InnerTree
}
