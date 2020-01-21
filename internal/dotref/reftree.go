package dotref

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/dotref/reftype"
)

type RefTree struct {
	cur    *RefTree
	parent *RefTree
	count  int
	nodes  []RefField
}

type RefField struct {
	FieldName string
	InnerTree *RefTree
	InnerRefs []RefField
	Index     int
	Type      reftype.RefType
}

func NewRefTree(parent *RefTree) *RefTree {
	newTree := RefTree{
		nodes:  []RefField{},
		parent: parent,
	}
	newTree.cur = &newTree
	return &newTree
}

func (r *RefTree) AddField(field string) error {
	rt := reftype.Unset
	if strings.HasPrefix(field, "$") {
		if r.cur.count > 0 {
			return fmt.Errorf("variable field is not allowed in dot-format reference: %s", field)
		}
		rt = reftype.VarName
	} else {
		rt = reftype.Literal
	}
	r.AddNode(RefField{FieldName: field, Type: rt})
	return nil
}

func (r *RefTree) AddNode(rf RefField) {
	r.cur.nodes = append(r.cur.nodes, rf)
	r.cur.count++
}

func (r *RefTree) CloseInner() error {
	if r.parent == nil {
		return fmt.Errorf("closing bracket without opeing one")
	}
	r.parent = r.parent.parent
	r.cur = r.parent
	return nil
}

func (r *RefTree) OpenInner() error {
	innterTree := NewRefTree(r)
	r.AddNode(RefField{InnerTree: innterTree, Type: reftype.InnerRef})
	// set cur to inner until bracket is closed
	r.parent = r.cur
	r.cur = innterTree
	return nil
}

func (r *RefTree) CreateResultFields() []RefField {
	for i := range r.nodes {
		if r.nodes[i].Type == reftype.InnerRef {
			r.nodes[i].InnerRefs = r.nodes[i].InnerTree.CreateResultFields()
			r.nodes[i].InnerTree = nil
		}
	}
	return r.nodes
}
