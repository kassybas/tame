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
	if len(fields) != 0 {
		return nil, fmt.Errorf("setting value of scalar variable with fields: %s %v", v.name, fields)
	}
	v.value = value
	return v, nil
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
