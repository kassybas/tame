package texpression

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/tame/internal/build/exprparse"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/types/exprtype"
)

type ExprField struct {
	FieldName string
	InnerRefs []ExprField
	Index     int
	Type      exprtype.ExprType
}

func NewField(val interface{}) (ExprField, error) {
	var newField ExprField
	switch val := val.(type) {
	case string:
		if strings.HasPrefix(val, "$") {
			// variable
			newField.Type = exprtype.VarName
			newField.FieldName = val
		} else if idx, err := strconv.Atoi(val); err == nil {
			// index
			newField.Type = exprtype.Index
			newField.Index = idx
		} else {
			// literal
			newField.Type = exprtype.Literal
			newField.FieldName, err = helpers.TrimLiteralQuotes(val)
			if err != nil {
				return newField, err
			}
		}
	case int:
		newField.Type = exprtype.Index
		newField.Index = val
	default:
		return newField, fmt.Errorf("unknown field type: %T", val)
	}
	return newField, nil
}

func NewExpression(expression string) ([]ExprField, error) {
	tree, err := exprparse.ParseExpression(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %s\n\t%s", expression, err.Error())
	}
	fields, err := newInnerExpression(tree)
	if err != nil {
		return fields, err
	}
	return fields, nil
}

func newInnerExpression(tree exprparse.ParseTree) ([]ExprField, error) {
	fields := []ExprField{}
	for i := range tree.Nodes {
		if tree.Nodes[i].InnerTree != nil {
			innerFields, err := newInnerExpression(*tree.Nodes[i].InnerTree)
			if err != nil {
				return nil, err
			}
			newField := ExprField{
				Type:      exprtype.InnerRef,
				InnerRefs: innerFields,
			}
			fields = append(fields, newField)
		} else {
			field, err := NewField(tree.Nodes[i].Val)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		}
	}
	return fields, nil
}