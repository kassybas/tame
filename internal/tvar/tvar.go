package tvar

import (
	"log"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/types/vartype"
)

type TVariable interface {
	Type() vartype.TVarType
	Name() string
	ToStr() string
	Value() interface{}
	ToEnvVars(string) []string
	SetValue([]dotref.RefField, interface{}) (TVariable, error)
}

func NewVariable(name string, value interface{}) TVariable {
	var newVar TVariable
	switch value.(type) {
	case []interface{}:
		{
			newVar = NewList(name, value.([]interface{}))
		}
	case bool, int, float64, string:
		{
			newVar = NewScalar(name, value)
		}
	case map[interface{}]interface{}, map[interface{}]TVariable, TVariable, []TVariable:
		{
			// encapsulate value to field
			newVar = NewMap(name, value)
		}
	default:
		{
			log.Fatalf("unknown variable type %s: %T", name, value)
		}

	}
	return newVar
}

func CopyVariable(newName string, sourceVar TVariable) TVariable {
	return NewVariable(newName, sourceVar.Value())
}
