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

// func CreateCompositeVariable(name string, value interface{}) TVariable {
// 	fields := strings.Split(name, keywords.TameFieldSeparator)
// 	last := len(fields) - 1
// 	innerVar := CreateVariable(fields[last], value)
// 	outerVar := CreateVariable(strings.Join(fields[:last], keywords.TameFieldSeparator), innerVar)
// 	return outerVar
// }

// func CreateListFromBracketsName(name string, value interface{}) (TVariable, error) {
// 	var tl TList
// 	index, n, err := helpers.ParseIndex(name)
// 	if err != nil {
// 		return tl, err
// 	}
// 	tl.name = n
// 	tl.value = make([]TVariable, index+1)
// 	for i := range tl.value {
// 		// Null all values other than the index
// 		tl.value[i] = TNull{}
// 	}
// 	tl.value[index] = CreateVariable(strconv.Itoa(index), value)
// 	return tl, nil
// }

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
