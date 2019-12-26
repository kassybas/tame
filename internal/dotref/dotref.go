package dotref

import (
	"fmt"
	"strconv"
	"strings"

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
			index, f, indexRef, err := parseIndex(field)
			if err != nil {
				return fields, err
			}
			fields = append(fields, RefField{FieldName: f})
			fields = append(fields, RefField{Index: index, FieldName: indexRef})
			continue
		}
		fields = append(fields, RefField{FieldName: field})
	}
	return fields, nil
}

// returns the index and variable name and index variable if exists
// if index is a variable reference, -1 is returned
func parseIndex(name string) (int, string, string, error) {
	lBr := strings.Index(name, keywords.IndexingSeparatorL) + 1
	rBr := strings.Index(name, keywords.IndexingSeparatorR)
	if strings.HasPrefix(name[lBr:rBr], keywords.PrefixReference) {
		// variable index
		return -1, name[0 : lBr-1], name[lBr:rBr], nil
	}
	index, err := strconv.Atoi(name[lBr:rBr])
	if err != nil {
		return 0, "", "", fmt.Errorf("not integer index or variable index: %s %s", name, name[lBr:rBr])
	}
	return index, name[0 : lBr-1], "", nil
}
