package tvar

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/texpression"
	"github.com/sirupsen/logrus"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TList struct {
	name   string
	values []TVariable
}

func NewList(name string, values interface{}) TList {
	tl := TList{
		name:   name,
		values: []TVariable{},
	}
	// this is repetitive because []interface{} is not a generic type enymore
	// it cannot be casted to []int like interface{}
	switch values := values.(type) {
	case []interface{}: //, []int, []string, []bool, []float64:
		for i := range values {
			tl.values = append(tl.values, NewVariable(ConvertKeyToString(i), values[i]))
		}
	case []int:
		for i := range values {
			tl.values = append(tl.values, NewVariable(ConvertKeyToString(i), values[i]))
		}
	case []string:
		for i := range values {
			tl.values = append(tl.values, NewVariable(ConvertKeyToString(i), values[i]))
		}
	case []float64:
		for i := range values {
			tl.values = append(tl.values, NewVariable(ConvertKeyToString(i), values[i]))
		}
	case []bool:
		for i := range values {
			tl.values = append(tl.values, NewVariable(ConvertKeyToString(i), values[i]))
		}
	default:
		logrus.Fatalf("internal error: creating list from non-iterable type: %s -- type: %T", name, values)
	}
	return tl
}

func (v TList) Name() string {
	return v.name
}

func (v TList) ToStr() string {
	return fmt.Sprintf("%v", v.values)
}

func (v TList) Value() interface{} {
	ifValues := []interface{}{}
	for i := range v.values {
		ifValues = append(ifValues, v.values[i].Value())
	}
	return ifValues
}

func (v TList) Type() vartype.TVarType {
	return vartype.TListType
}

func (v TList) GetItem(i int) (interface{}, error) {
	if len(v.values) <= i {
		return nil, fmt.Errorf("index out of range: %s[%d]", v.name, i)
	}
	return v.values[i], nil
}

func (v TList) ToEnvVars(ShellFieldSeparator string) []string {
	var envVars []string
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	for _, item := range v.values {
		// TODO: handle brackets better
		for _, itemEnvVar := range item.ToEnvVars(ShellFieldSeparator) {
			ev := fmt.Sprintf("%s%s", trimmedName, itemEnvVar)
			envVars = append(envVars, ev)
		}
	}
	return envVars
}

func (v TList) SetValue(fields []texpression.ExprField, value interface{}) (TVariable, error) {
	var err error
	if len(fields) == 0 {
		// this should never happen, since this would mean that texpression field was called with empty string
		// which is an invalid variable name, which fails at parsing
		return nil, fmt.Errorf("internal error: empty reference")
	}
	if len(fields) == 1 {
		// overdefine list with different type
		// if no fieldname??
		return NewVariable(fields[0].FieldName, value), nil
	}
	field := fields[1]
	if field.FieldName != "" {
		return nil, fmt.Errorf("referencing field on a list: %s %v ", v.name, fields)
	}
	if field.Index >= len(v.values) || field.Index < 0 {
		return nil, fmt.Errorf("index out-of-range: %s[%d]", v.name, field.Index)
	}
	v.values[field.Index], err = v.values[field.Index].SetValue(fields[1:], value)
	return v, err
}

func (v TList) GetInnerValue(fields []texpression.ExprField) (interface{}, error) {
	if len(fields) == 0 {
		// this should never happen, since this would mean that texpression field was called with empty string
		// which is an invalid variable name, which fails at parsing
		return nil, fmt.Errorf("internal error: empty reference")
	}
	if len(fields) == 1 {
		// field[0] is the variable name
		return v.Value(), nil
	}
	// field[1] is the first actual field, should be the index
	field := fields[1]
	if field.FieldName != "" {
		return nil, fmt.Errorf("referencing field on a list: %s %v ", v.name, fields)
	}
	if field.Index >= len(v.values) || field.Index < 0 {
		return nil, fmt.Errorf("index out-of-range: %s[%d]", v.name, field.Index)
	}
	return v.values[field.Index].GetInnerValue(fields[1:])
}
