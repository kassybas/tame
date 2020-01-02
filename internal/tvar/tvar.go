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
	GetInnerValue([]dotref.RefField) (interface{}, error)
}

func NewVariable(name string, value interface{}) TVariable {
	var newVar TVariable
	switch value := value.(type) {
	case []interface{}:
		{
			newVar = NewList(name, value)
		}
	case bool, int, float64, string, nil:
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
