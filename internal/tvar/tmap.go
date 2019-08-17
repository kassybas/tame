package tvar

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
)

type TMap struct {
	name  string
	value map[string]VariableI
}

func (v TMap) Type() TVarType {
	return TMapType
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

func CreateMap(name string, value map[interface{}]interface{}) TMap {
	var tm TMap
	tm.name = name
	tm.value = make(map[string]VariableI)
	for k, v := range value {
		tm.value[k.(string)] = CreateVariable(k.(string), v)
	}
	return tm
}

func UpdateCompositeValue(origVar, newField VariableI) VariableI {
	if newField.Type() != TMapType {
		return newField
	}
	origM := origVar.(TMap)
	newM := newField.(TMap)

	for k, newVal := range newM.value {
		member, exists := origM.value[k]
		if exists && member.Type() != TMapType {
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
		if newList.value[i].Type() != TNullType {
			origList.value[i] = newList.value[i]
		}
	}
	return origList, nil
}
