package tvar

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/dotref"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TList struct {
	name   string
	values []TVariable
}

func NewList(name string, values []interface{}) TList {
	tl := TList{
		name:   name,
		values: make([]TVariable, len(values)),
	}
	for i := range values {
		tl.values[i] = NewVariable(ConvertKeyToString(i), values[i])
	}
	return tl
}

func (v TList) Name() string {
	return v.name
}

func (v TList) ToStr() string {
	// TODO: yaml dump?
	return fmt.Sprintf("%v", v.values)
}

func (v TList) Value() interface{} {
	return v.values
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

func (v TList) SetValue(fields []dotref.RefField, value interface{}) (TVariable, error) {
	var err error
	field := fields[0]
	if field.FieldName != "" {
		return nil, fmt.Errorf("setting map on a list: [%s] %v ", v.name, fields)
	}
	if field.Index >= len(v.values) {
		return nil, fmt.Errorf("index out-of-range: %s[%d]", v.name, field.Index)
	}
	v.values[field.Index], err = v.values[field.Index].SetValue(fields[1:], value)
	return v, err
}
