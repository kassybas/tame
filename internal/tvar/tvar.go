package tvar

import (
	"github.com/sirupsen/logrus"
)

type VariableI interface {
	Type() TVarType
	Name() string
	Value() interface{}
	ToInt() (int, error)
	ToStr() string
	ToEnvVars() []string
	IsScalar() bool
}

// type Variable struct {
// 	Name string
// 	// TODO: interface
// 	stringValue string
// 	intValue    int
// 	Type        TVarType
// }

func CreateVariable(name string, value interface{}) VariableI {
	switch value.(type) {
	// Null
	case nil:
		{
			return TNull{name: name}
		}
	case TNull:
		{
			return TNull{name: name}
		}
	// Bool
	case bool:
		{
			return TBool{name: name, value: value.(bool)}
		}
	case TBool:
		{
			tb := value.(TBool)
			return TBool{name: name, value: tb.value}
		}
	// String
	case string:
		{
			return TString{name: name, value: value.(string)}
		}
	case TString:
		{
			ts := value.(TString)
			return TString{name: name, value: ts.value}
		}
	// Int
	case int:
		{
			return TInt{name: name, value: value.(int)}
		}
	case TInt:
		{
			ti := value.(TInt)
			return TInt{name: name, value: ti.value}
		}
		// Float
	case float64:
		{
			return TFloat{name: name, value: value.(float64)}
		}
	case TFloat:
		{
			tf := value.(TFloat)
			return TFloat{name: name, value: tf.value}
		}
	// Map
	case map[interface{}]interface{}:
		{
			return CreateMap(name, value.(map[interface{}]interface{}))
		}
	case TMap:
		{
			tm := value.(TMap)
			return TMap{name: name, value: tm.value}
		}
	// List
	case []interface{}:
		{
			return CreateList(name, value.([]interface{}))
		}
	case TList:
		{
			tl := value.(TList)
			return TList{name: name, value: tl.value}
		}
	default:
		{
			logrus.Fatalf("Undeterminable variable type for: %s -- %T", name, value)
		}
	}
	return nil
}

type TVarType int

const (
	TErrorType TVarType = iota
	TStringType
	TIntType
	TFloatType
	TListType
	TMapType
	TBoolType
	TNullType
)
