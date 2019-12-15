package tvar

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TMap struct {
	name  string
	value map[string]VariableI
}

func (v TMap) Type() vartype.TVarType {
	return vartype.TMapType
}

func (v TMap) IsScalar() bool {
	return false
}

func (v TMap) Name() string {
	return v.name
}

func (v TMap) Value() interface{} {
	return v.value
}

func (v TMap) ToInt() (int, error) {
	// i, err := strconv.Atoi(v.value.(string))
	// return i, err
	return 0, nil
}

func (v TMap) ToStr() string {
	return ""
	// return v.value.(string)
}

func (v TMap) ToEnvVars() []string {
	var envVars []string
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	for _, v := range v.value {
		for _, memberEnvVar := range v.ToEnvVars() {
			ev := trimmedName + keywords.ShellFieldSeparator + memberEnvVar
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

func CreateMap(name string, value map[interface{}]interface{}) TMap {
	var tm TMap
	tm.name = name
	tm.value = make(map[string]VariableI)
	for k, v := range value {
		tm.value[k.(string)] = CreateVariable(k.(string), v)
	}
	return tm
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
	return TMap{name: name, value: map[string]VariableI{innerValue.Name(): innerValue}}
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
