package tvar

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TMap struct {
	TBaseVar
	value map[string]VariableI
}

func convertKeyToString(key interface{}) string {
	switch key.(type) {
	case string:
		{
			return key.(string)
		}
	case int:
		{
			return strconv.Itoa(key.(int))
		}
	case float64:
		{
			return fmt.Sprintf("%g", key.(float64))
		}
	case bool:
		{
			if key.(bool) {
				return "true"
			}
			return "false"
		}
	}
	log.Fatal("yaml key not valid type (must be: string, int, float, bool)", key)
	return ""
}

func NewMap(name string, value map[interface{}]interface{}) TMap {
	tm := TMap{
		TBaseVar: TBaseVar{
			name:    name,
			iValue:  value,
			varType: vartype.TMapType,
		},
	}
	tm.value = make(map[string]VariableI)
	for k, v := range value {
		// original type of key is converted to string
		// TODO: make this more flexible
		stringKey := convertKeyToString(k)
		tm.value[stringKey] = CreateVariable(stringKey, v)
	}
	return tm
}

func (v TMap) ToInt() (int, error) {
	return 0, nil
}

func (v TMap) ToStr() string {
	return ""
	// return v.value.(string)
}

func (v TMap) ToEnvVars(ShellFieldSeparator string) []string {
	var envVars []string
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	for _, v := range v.value {
		for _, memberEnvVar := range v.ToEnvVars(ShellFieldSeparator) {
			ev := fmt.Sprintf("%s%s%s", trimmedName, ShellFieldSeparator, memberEnvVar)
			envVars = append(envVars, ev)
		}
	}
	return envVars
}

func (tm TMap) IsMember(key string) bool {
	_, exist := tm.value[key]
	return exist
}

func (tm TMap) GetMember(key string) (VariableI, error) {
	v, exist := tm.value[key]
	if !exist {
		return nil, fmt.Errorf("field is not member of map: %s.%s", tm.name, key)
	}
	return v, nil
}

func ValidateUpdate(origVar VariableI, dr dotref.DotRef) error {
	if len(dr.Fields) == 0 {
		return nil
	}
	var err error
	cur := origVar
	lastField := len(dr.Fields) - 1
	for i, field := range dr.Fields {
		if field.FieldName == "" {
			if cur.Type() != vartype.TListType {
				return fmt.Errorf("indexing non-list type: %s[%d] (type: %s)", cur.Name(), field.Index, vartype.GetTypeNameString(cur.Type()))
			}
			cur, err = cur.(TList).GetItem(field.Index)
			if err != nil {
				return err
			}
			continue
		}
		if cur.Type() != vartype.TMapType {
			return fmt.Errorf("field reference on a non-map type: %s.%s (type: %s)", cur.Name(), field.FieldName, vartype.GetTypeNameString(cur.Type()))
		}
		// last fields can be added to the map
		if i != lastField {
			cur, err = cur.(TMap).GetMember(field.FieldName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateCompositeValue(origVar, newField VariableI) VariableI {
	if newField.Type() != vartype.TMapType {
		return newField
	}
	origM := origVar.(TMap)
	newM := newField.(TMap)

	for k, newVal := range newM.value {
		member, exists := origM.value[k]
		if exists && member.Type() != vartype.TMapType {
			// scalar update
			origM.value[k] = newVal
			continue
		}
		if exists {
			// merge with the member value
			origM.value[k] = UpdateCompositeValue(member, newVal)
			continue
		}
		origM.value[k] = newVal
	}
	return origM
}

func EncapsulateValueToMap(name string, innerValue VariableI) TMap {
	m := make(map[interface{}]interface{})
	m[innerValue.Name()] = innerValue.Value()
	return NewMap(name, m)
}

func MergeLists(origList, newList TList) (TList, error) {
	if len(origList.value) < len(newList.value) {
		return TList{}, fmt.Errorf("index out of range: %s[%d]", origList.name, len(newList.value)-1)
	}
	for i := range newList.value {
		if newList.value[i].Type() != vartype.TNullType {
			origList.value[i] = newList.value[i]
		}
	}
	return origList, nil
}
