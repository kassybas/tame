package dotref

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/types/reftype"
)

type RefTreeParse struct {
	cur    *RefTreeParse
	parent *RefTreeParse
	count  int
	nodes  []RefField
}

func NewRefTree(parent *RefTreeParse) *RefTreeParse {
	newTree := RefTreeParse{
		nodes:  []RefField{},
		parent: parent,
	}
	newTree.cur = &newTree
	return &newTree
}

func trimLiteralQuotes(field string) (string, error) {
	if strings.HasPrefix(field, `"`) {
		if !strings.HasSuffix(field, `"`) {
			return "", fmt.Errorf("missing closing bracket: %s", field)
		}
		field = strings.Trim(field, `"`)
	} else if strings.HasPrefix(field, `'`) {
		if !strings.HasSuffix(field, `"`) {
			return "", fmt.Errorf("missing closing bracket: %s", field)
		}
		field = strings.Trim(field, `'`)
	}
	return field, nil
}

func (r *RefTreeParse) AddField(field string) error {
	newField, err := NewField(field)
	if err != nil {
		return err
	}
	if newField.Type == reftype.VarName && r.cur.count > 0 {
		return fmt.Errorf("variable field is not allowed in dot-format reference: %s", field)
	}
	r.AddNode(newField)
	return nil
}

func (r *RefTreeParse) AddNode(rf RefField) {
	r.cur.nodes = append(r.cur.nodes, rf)
	r.cur.count++
}

func (r *RefTreeParse) CloseInner() error {
	if r.parent == nil {
		return fmt.Errorf("closing bracket without opeing one")
	}
	r.parent = r.parent.parent
	r.cur = r.parent
	return nil
}

func (r *RefTreeParse) OpenInner() error {
	innterTree := NewRefTree(r)
	r.AddNode(RefField{InnerTree: innterTree, Type: reftype.InnerRef})
	// set cur to inner until bracket is closed
	r.parent = r.cur
	r.cur = innterTree
	return nil
}

func (r *RefTreeParse) CreateResultFields() []RefField {
	for i := range r.nodes {
		if r.nodes[i].Type == reftype.InnerRef {
			r.nodes[i].InnerRefs = r.nodes[i].InnerTree.CreateResultFields()
			r.nodes[i].InnerTree = nil
		}
	}
	return r.nodes
}
