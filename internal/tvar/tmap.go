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
	if len(fields) < 2 {
		return nil, fmt.Errorf("setting scalar on a map: [%s] %v ", v.name, fields)
	}
	field := fields[0]
	if field.FieldName == "" {
		return nil, fmt.Errorf("empty field name: %s%s", v.name, keywords.TameFieldSeparator)
	}
	// Map field
	val, exists := v.values[field]
	if !exists {
		return nil, fmt.Errorf("field does not exist in map: %s%s%s ", v.name, keywords.TameFieldSeparator, field)
	}
	v.values[field.FieldName], err = val.SetValue(fields[1:], value)
	return v, err
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

// func (tm TMap) GetMember(key string) (TVariable, error) {
// 	v, exist := tm.value[key]
// 	if !exist {
// 		return nil, fmt.Errorf("field is not member of map: %s.%s", tm.name, key)
// 	}
// 	return v, nil
// }

// TODO: cleanup
// func ValidateUpdate(origVar TVariable, dr dotref.DotRef) error {
// 	if len(dr.Fields) == 0 {
// 		return nil
// 	}
// 	var err error
// 	cur := origVar
// 	lastField := len(dr.Fields) - 1
// 	for i, field := range dr.Fields {
// 		if field.FieldName == "" {
// 			if cur.Type() != vartype.TListType {
// 				return fmt.Errorf("indexing non-list type: %s[%d] (type: %s)", cur.Name(), field.Index, vartype.GetTypeNameString(cur.Type()))
// 			}
// 			cur, err = cur.(TList).GetItem(field.Index)
// 			if err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		if cur.Type() != vartype.TMapType {
// 			return fmt.Errorf("field reference on a non-map type: %s.%s (type: %s)", cur.Name(), field.FieldName, vartype.GetTypeNameString(cur.Type()))
// 		}
// 		// last fields can be added to the map
// 		if i != lastField {
// 			cur, err = cur.(TMap).GetMember(field.FieldName)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func UpdateCompositeValue(origVar, newField TVariable) TVariable {
// 	if newField.Type() != vartype.TMapType {
// 		return newField
// 	}
// 	origM := origVar.(TMap)
// 	newM := newField.(TMap)

// 	for k, newVal := range newM.value {
// 		member, exists := origM.value[k]
// 		if exists && member.Type() != vartype.TMapType {
// 			// scalar update
// 			origM.value[k] = newVal
// 			continue
// 		}
// 		if exists {
// 			// merge with the member value
// 			origM.value[k] = UpdateCompositeValue(member, newVal)
// 			continue
// 		}
// 		origM.value[k] = newVal
// 	}
// 	return origM
// }

// func EncapsulateValueToMap(name string, innerValue TVariable) TMap {
// 	m := make(map[interface{}]interface{})
// 	m[innerValue.Name()] = innerValue.Value()
// 	return NewMap(name, m)
// }
