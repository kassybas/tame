package dotref

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/helpers"

	"github.com/kassybas/tame/internal/keywords"
)

type RefField struct {
	FieldName string
	Index     int
}

func Empty() []RefField {
	return []RefField{}
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

func ParseFields(dotName string) ([]RefField, error) {
	tmp_fields := strings.Split(dotName, keywords.TameFieldSeparator)
	fields := []RefField{}

	for _, field := range tmp_fields {
		hasIndex, err := checkIndexSeparators(field)
		if err != nil {
			return fields, err
		}
		if hasIndex {
			index, f, err := helpers.ParseIndex(field)
			if err != nil {
				return fields, err
			}
			fields = append(fields, RefField{FieldName: f})
			fields = append(fields, RefField{Index: index})
		} else {
			fields = append(fields, RefField{FieldName: field})
		}
	}
	return fields, nil
}
