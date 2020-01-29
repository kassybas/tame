package dotref

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/tame/types/reftype"
)

type RefField struct {
	// FieldName interface{}
	FieldName string
	InnerTree *RefTreeParse
	InnerRefs []RefField
	Index     int
	Type      reftype.RefType
}

func NewField(val interface{}) (RefField, error) {
	var newField RefField
	switch val := val.(type) {
	case string:
		if strings.HasPrefix(val, "$") {
			// variable
			newField.Type = reftype.VarName
			newField.FieldName = val
		} else if idx, err := strconv.Atoi(val); err == nil {
			// index
			newField.Type = reftype.Index
			newField.Index = idx
		} else {
			// literal
			newField.Type = reftype.Literal
			newField.FieldName, err = trimLiteralQuotes(val)
			if err != nil {
				return newField, err
			}
		}
	case int:
		newField.Type = reftype.Index
		newField.Index = val
	default:
		return newField, fmt.Errorf("unknown field type: %T", val)
	}
	return newField, nil
}
