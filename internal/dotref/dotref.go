package dotref

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/helpers"

	"github.com/kassybas/mate/internal/keywords"
)

type Field struct {
	FieldName string
	Index     int
}

type DotRef struct {
	Fields []Field
	Name   string
	Value  interface{}
}

func checkIndexSeparators(field string) (bool, error) {
	if !strings.Contains(field, keywords.IndexingSeparatorL) && !strings.Contains(field, keywords.IndexingSeparatorR) {
		// no indexing in field
		return false, nil
	}
	if strings.Contains(field, keywords.IndexingSeparatorL) != strings.Contains(field, keywords.IndexingSeparatorR) {
		// one of them is not in the string
		return false, fmt.Errorf("missing bracket: %s", field)
	}
	return true, nil
}

func IsDotRef(dotName string) bool {
	if strings.Contains(dotName, keywords.IndexingSeparatorL) || strings.Contains(dotName, keywords.IndexingSeparatorR) ||
		strings.Contains(dotName, keywords.TameFieldSeparator) {
		return true
	}
	return false

}

func NewReference(dotName string, value interface{}) (DotRef, error) {
	var dv DotRef
	// if !IsDotRef(dotName) {
	// 	return dv, fmt.Errorf("internal error: creatign reference from non-dot name %s", dotName)
	// }
	fields := strings.Split(dotName, keywords.TameFieldSeparator)

	firstField := fields[0]
	dv.Name = strings.Split(firstField, keywords.IndexingSeparatorL)[0]

	dv.Fields = []Field{}
	for i, field := range fields {
		hasIndex, err := checkIndexSeparators(field)
		if err != nil {
			return dv, err
		}
		if hasIndex {
			index, f, err := helpers.ParseIndex(field)
			if err != nil {
				return dv, err
			}
			// skipping first field because that is the name
			if i != 0 {
				dv.Fields = append(dv.Fields, Field{FieldName: f})
			}
			dv.Fields = append(dv.Fields, Field{Index: index})
		} else if i != 0 {
			f := field
			dv.Fields = append(dv.Fields, Field{FieldName: f})
		}
	}
	dv.Value = value
	return dv, nil
}
