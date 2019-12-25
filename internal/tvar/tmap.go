package tvar

import (
	"fmt"
	"log"
	"strings"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TMap struct {
	name   string
	values map[interface{}]TVariable
}

func ConvertKeyToString(key interface{}) string {
	// TODO: using map or list as a key can be difficult to figure out
	//       refering to it via the "%v" in a dotref is cumbersome, fix that
	return fmt.Sprintf("%v", key)
}

func NewMap(name string, value interface{}) TMap {
	tm := TMap{
		name:   name,
		values: make(map[interface{}]TVariable),
	}
	switch value.(type) {
	case map[interface{}]interface{}:
		{
			// create map
			for k, v := range value.(map[interface{}]interface{}) {
				// original type of key is converted to string in the name
				stringKey := ConvertKeyToString(k)
				tm.values[k] = NewVariable(stringKey, v)
			}
		}
	case map[interface{}]TVariable:
		{
			tm.values = value.(map[interface{}]TVariable)
		}
	case TVariable:
		{
			// encapsulate var to map
			innerVar := value.(TVariable)
			tm.values[innerVar.Name()] = tm

		}
	default:
		{
			log.Fatalf("unknown type to create map from %s: %T", name, value)
		}
	}
	return tm
}

func (v TMap) Name() string {
	return v.name
}

func (v TMap) ToStr() string {
	// TODO: yaml dump?
	return fmt.Sprintf("%v", v.values)
}

func (v TMap) Type() vartype.TVarType {
	return vartype.TScalarType
}

func (v TMap) Value() interface{} {
	return v.values
}

func (v TMap) SetValue(fields []dotref.RefField, value interface{}) (TVariable, error) {
	var err error
	if len(fields) == 0 {
		// this should never happen, since this would mean that dotref field was called with empty string
		return nil, fmt.Errorf("internal error: getting member with a non-reference field")
	}
	if len(fields) == 1 {
		// cast value of map to different value
		return NewVariable(fields[0].FieldName, value), nil
	}
	// Field 0 is the name of the variable
	field := fields[1]
	if field.FieldName == "" {
		return nil, fmt.Errorf("empty field name -- indexing not allowed on map type: %s index: %d", v.name, field.Index)
	}
	// Map field
	if !v.IsMember(field.FieldName) {
		if len(fields) == 2 {
			// last field's key's can be extended
			v.values[fields[1].FieldName] = NewVariable(fields[1].FieldName, value)
			return v, nil
		}
		return nil, fmt.Errorf("field does not exist in map: %s: %s.%s ", v.name, fields[0].FieldName, fields[1].FieldName)
	}
	// Setting an existing member
	v.values[field.FieldName], err = v.values[field.FieldName].SetValue(fields[1:], value)
	return v, err
}

func (v TMap) GetInnerValue(fields []dotref.RefField) (interface{}, error) {
	if len(fields) == 0 {
		// this should never happen, since this would mean that dotref field was called empty
		return nil, fmt.Errorf("internal error: empty reference")
	}
	if len(fields) == 1 {
		// field[0] is the variable name
		return v.Value(), nil
	}
	// field[1] is the first actual field
	field := fields[1]
	if field.FieldName == "" {
		return nil, fmt.Errorf("empty field name -- indexing not allowed on map type: %s index: %d", v.name, field.Index)
	}
	if !v.IsMember(field.FieldName) {
		return nil, fmt.Errorf("field does not exist %s.%s", v.name, field.FieldName)
	}
	return v.values[field.FieldName].GetInnerValue(fields[1:])
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
