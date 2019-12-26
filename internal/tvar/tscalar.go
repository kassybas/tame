package tvar

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/dotref"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TScalar struct {
	name    string
	value   interface{}
	varType vartype.TVarType
}

func NewScalar(name string, value interface{}) TScalar {
	return TScalar{
		name:  name,
		value: value,
	}
}

func (v TScalar) Name() string {
	return v.name
}

func (v TScalar) Value() interface{} {
	return v.value
}

func (v TScalar) SetValue(fields []dotref.RefField, value interface{}) (TVariable, error) {
	if len(fields) == 0 {
		// this should never happen, since this would mean that dotref field was called with empty string
		// which is an invalid variable name, which fails at parsing
		return nil, fmt.Errorf("internal error: empty reference")
	}
	if len(fields) > 1 {
		return nil, fmt.Errorf("field reference on scalar variable (only allowed on map or list): '%s'.'%v'", v.name, fields[1].FieldName)
	}
	v.value = value
	return v, nil
}

func (v TScalar) GetInnerValue(fields []dotref.RefField) (interface{}, error) {
	if len(fields) == 0 {
		// this should never happen, since this would mean that dotref field was called with empty string
		// which is an invalid variable name, which fails at parsing
		return nil, fmt.Errorf("internal error: empty reference")
	}
	if len(fields) != 1 {
		return nil, fmt.Errorf("field reference on scalar variable (only allowed on map or list): '%s'.'%v'", v.name, fields[1].FieldName)
	}
	return v.value, nil
}

func (v TScalar) ToStr() string {
	return fmt.Sprintf("%v", v.value)
}

func (v TScalar) Type() vartype.TVarType {
	return vartype.TScalarType
}

func (v TScalar) ToEnvVars(ShellFieldSeparator string) []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	if v.value == nil {
		return []string{}
	}
	return []string{
		fmt.Sprintf("%s=%s", trimmedName, v.ToStr()),
	}
}
