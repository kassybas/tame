package tvar

import (
	"strings"

	"github.com/kassybas/mate/internal/keywords"
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

func CreateCompositeVariable(name string, value interface{}) VariableI {
	fields := strings.Split(name, keywords.TameFieldSeparator)
	last := len(fields) - 1
	innerVar := CreateVariable(fields[last], value)
	outerVar := CreateVariable(strings.Join(fields[:last], keywords.TameFieldSeparator), innerVar)
	return outerVar
}

func CreateVariable(name string, value interface{}) VariableI {
	if strings.Contains(name, keywords.TameFieldSeparator) || strings.Contains(name, keywords.IndexingSeparatorL) || strings.Contains(name, keywords.IndexingSeparatorR) {
		return CreateCompositeVariable(name, value)
	}
	switch value.(type) {
	// Null
	case nil:
		{
			return TNull{name: name}
		}
	case TNull:
		{
			return EncapsulateValueToMap(name, value.(TNull))
		}
	// Bool
	case bool:
		{
			return TBool{name: name, value: value.(bool)}
		}
	case TBool:
		{
			tb := value.(TBool)
			return EncapsulateValueToMap(name, tb)
		}
	// String
	case string:
		{
			return TString{name: name, value: value.(string)}
		}
	case TString:
		{
			ts := value.(TString)
			return EncapsulateValueToMap(name, ts)
		}
	// Int
	case int:
		{
			return &TInt{name: name, value: value.(int)}
		}
	case TInt:
		{
			ti := value.(TInt)
			return EncapsulateValueToMap(name, ti)
		}
	// Float
	case float64:
		{
			return TFloat{name: name, value: value.(float64)}
		}
	case TFloat:
		{
			tf := value.(TFloat)
			return EncapsulateValueToMap(name, tf)
		}
	// Map
	case map[interface{}]interface{}:
		{
			m := CreateMap(name, value.(map[interface{}]interface{}))
			return m
		}
	case map[string]VariableI:
		{
			return TMap{name: name, value: value.(map[string]VariableI)}
		}
	case TMap:
		{
			tm := value.(TMap)
			return EncapsulateValueToMap(name, tm)
		}
	// List
	case []interface{}:
		{
			l := CreateList(name, value.([]interface{}))
			return l
		}
	case TList:
		{
			tl := value.(TList)
			return EncapsulateValueToMap(name, tl)
		}
	default:
		{
			logrus.Fatalf("Undeterminable variable type for: %s -- %T", name, value)
		}
	}
	return nil
}

func CopyVariable(newName string, sourceVar VariableI) VariableI {
	return CreateVariable(newName, sourceVar.Value())
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
