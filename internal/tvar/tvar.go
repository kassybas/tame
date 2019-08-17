package tvar

import (
	"fmt"
	"strconv"
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

// returns the index and variable name
func parseIndex(name string) (int, string, error) {
	lBr := strings.Index(name, keywords.IndexingSeparatorL) + 1
	rBr := strings.Index(name, keywords.IndexingSeparatorR)
	index, err := strconv.Atoi(name[lBr:rBr])
	if err != nil {
		return 0, "", fmt.Errorf("not integer index: %s %s", name, name[lBr:rBr])
	}
	return index, name[0 : lBr-1], nil
}

func CreateListFromBracketsName(name string, value interface{}) (VariableI, error) {
	var tl TList
	index, n, err := parseIndex(name)
	if err != nil {
		return tl, err
	}
	tl.name = n
	tl.value = make([]VariableI, index+1)
	for i := range tl.value {
		// Null all values other than the index
		tl.value[i] = TNull{}
	}
	tl.value[index] = CreateVariable(strconv.Itoa(index), value)
	return tl, nil
}

func CreateVariable(name string, value interface{}) VariableI {
	if strings.Contains(name, keywords.TameFieldSeparator) {
		return CreateCompositeVariable(name, value)
	}
	if strings.Contains(name, keywords.IndexingSeparatorL) && strings.Contains(name, keywords.IndexingSeparatorR) {
		l, err := CreateListFromBracketsName(name, value)
		if err != nil {
			logrus.Fatal(err)
		}
		return l
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
			return EncapsulateValueToMap(name, value.(TBool))
		}
	// String
	case string:
		{
			return TString{name: name, value: value.(string)}
		}
	case TString:
		{
			return EncapsulateValueToMap(name, value.(TString))
		}
	// Int
	case int:
		{
			return TInt{name: name, value: value.(int)}
		}
	case TInt:
		{
			return EncapsulateValueToMap(name, value.(TInt))
		}
	// Float
	case float64:
		{
			return TFloat{name: name, value: value.(float64)}
		}
	case TFloat:
		{
			return EncapsulateValueToMap(name, value.(TFloat))
		}
	// Map
	case map[interface{}]interface{}:
		{
			return CreateMap(name, value.(map[interface{}]interface{}))
		}
	case map[string]VariableI:
		{
			return TMap{name: name, value: value.(map[string]VariableI)}
		}
	case TMap:
		{
			return EncapsulateValueToMap(name, value.(TMap))
		}
	// List
	case []interface{}:
		{
			return CreateListFromInterface(name, value.([]interface{}))
		}
	case []VariableI:
		{
			return CreateListFromVars(name, value.([]VariableI))
		}
	case TList:
		{
			return EncapsulateValueToMap(name, value.(TList))
		}
	// Undefined
	default:
		{
			logrus.Fatalf("Internal error: undeterminable variable type for %s -- %T", name, value)
		}
	}
	return nil
}

func CopyVariable(newName string, sourceVar VariableI) VariableI {
	return CreateVariable(newName, sourceVar.Value())
}

type TVarType int

const (
	TUnknownType TVarType = iota
	TStringType
	TIntType
	TFloatType
	TListType
	TMapType
	TBoolType
	TNullType
)
