package tvar

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/texpression"
	"github.com/kassybas/tame/types/vartype"
	"github.com/sirupsen/logrus"
)

type TMap struct {
	name   string
	values map[interface{}]TVariable
}

func ConvertKeyToString(key interface{}) string {
	// TODO: using map or list as a key can be difficult to figure out
	//       refering to it via the "%v" in a texpression is cumbersome, fix that
	return fmt.Sprintf("%v", key)
}

func NewMap(name string, value interface{}) TMap {
	tm := TMap{
		name:   name,
		values: make(map[interface{}]TVariable),
	}
	switch value := value.(type) {
	case map[interface{}]interface{}:
		{
			// create map
			for k, v := range value {
				// original type of key is converted to string in the name of the variable not in the map
				stringKey := ConvertKeyToString(k)
				tm.values[k] = NewVariable(stringKey, v)
			}
		}
	case map[string]interface{}:
		{
			for k, v := range value {
				tm.values[k] = NewVariable(k, v)
			}
		}
	case map[interface{}]TVariable:
		{
			tm.values = value
		}
	case TVariable:
		{
			// encapsulate var to map
			innerVar := value
			tm.values[innerVar.Name()] = tm

		}
	default:
		{
			logrus.Fatalf("unknown type to create map from %s: %T", name, value)
		}
	}
	return tm
}

func (v TMap) Name() string {
	return v.name
}

func (v TMap) ToStr() string {
	return fmt.Sprintf("%v", v.values)
}

func (v TMap) Type() vartype.TVarType {
	return vartype.TMapType
}

func (v TMap) Value() interface{} {
	ifValues := make(map[interface{}]interface{})
	for k, v := range v.values {
		ifValues[k] = v.Value()
	}
	return ifValues
}

func (v TMap) SetValue(fields []texpression.ExprField, value interface{}) (TVariable, error) {
	var err error
	if len(fields) == 0 {
		// this should never happen, since this would mean that texpression field was called with empty string
		// which is an invalid variable name, which fails at parsing
		return nil, fmt.Errorf("internal error: getting member with a non-reference field")
	}
	if len(fields) == 1 {
		// cast value of map to different value
		return NewVariable(fields[0].Val, value), nil
	}
	// Field 0 is the name of the variable
	field := fields[1]
	if field.Val == "" {
		return nil, fmt.Errorf("empty field name -- indexing not allowed on map type: %s index: %d", v.name, field.Index)
	}
	// Map field
	if !v.IsMember(field.Val) {
		if len(fields) == 2 {
			// last field's key's can be extended
			v.values[fields[1].Val] = NewVariable(fields[1].Val, value)
			return v, nil
		}
		return nil, fmt.Errorf("field does not exist in map: %s: %s.%s ", v.name, fields[0].Val, fields[1].Val)
	}
	// Setting an existing member
	v.values[field.Val], err = v.values[field.Val].SetValue(fields[1:], value)
	return v, err
}

func (v TMap) GetInnerValue(fields []texpression.ExprField) (interface{}, error) {
	if len(fields) == 0 {
		// this should never happen, since this would mean that texpression field was called with empty string
		// which is an invalid variable name, which fails at parsing
		return nil, fmt.Errorf("internal error: empty reference")
	}
	if len(fields) == 1 {
		// field[0] is the variable name
		return v.Value(), nil
	}
	// field[1] is the first actual field
	field := fields[1]
	if field.Val == "" {
		return nil, fmt.Errorf("empty field name -- indexing not allowed on map type: %s index: %d", v.name, field.Index)
	}
	if !v.IsMember(field.Val) {
		return nil, fmt.Errorf("field does not exist %s.%s", v.name, field.Val)
	}
	return v.values[field.Val].GetInnerValue(fields[1:])
}

func (v TMap) ToEnvVars(ShellFieldSeparator string) []string {
	var envVars []string
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	for _, v := range v.values {
		for _, memberEnvVar := range v.ToEnvVars(ShellFieldSeparator) {
			ev := fmt.Sprintf("%s%s%s", trimmedName, ShellFieldSeparator, memberEnvVar)
			envVars = append(envVars, ev)
		}
	}
	return envVars
}

func (tm TMap) IsMember(key interface{}) bool {
	_, exist := tm.values[key]
	return exist
}
