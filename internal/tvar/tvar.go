package tvar

import (
	"github.com/kassybas/tame/internal/texpression"
	"github.com/kassybas/tame/types/vartype"
	"github.com/sirupsen/logrus"
)

type TVariable interface {
	Type() vartype.TVarType
	Name() string
	ToStr() string
	Value() interface{}
	ToEnvVars(string) []string
	SetValue([]texpression.ExprField, interface{}) (TVariable, error)
	GetInnerValue([]texpression.ExprField) (interface{}, error)
}

func NewVariable(name string, value interface{}) TVariable {
	var newVar TVariable
	switch value := value.(type) {
	case []interface{}, []int, []string, []float64:
		{
			newVar = NewList(name, value)
		}
	case bool, int, float64, string, nil:
		{
			newVar = NewScalar(name, value)
		}
	case map[interface{}]interface{}, map[string]interface{}, map[interface{}]TVariable, TVariable, []TVariable:
		{
			// encapsulate value to field
			newVar = NewMap(name, value)
		}
	default:
		{
			logrus.Fatalf("unknown variable type %s: %T", name, value)
		}

	}
	return newVar
}
