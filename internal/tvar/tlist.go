package tvar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TList struct {
	TBaseVar
	value []VariableI
}

func NewList(name string, value []VariableI) TList {
	return TList{
		TBaseVar: TBaseVar{
			name:     name,
			iValue:   interface{}(value),
			isScalar: false,
			varType:  vartype.TListType,
		},
		value: value,
	}

}

func (v TList) ToInt() (int, error) {
	// i, err := strconv.Atoi(v.value.(string))
	// return i, err
	return 0, nil
}

func (v TList) ToStr() string {
	return ""
}

func (v TList) GetItem(i int) (VariableI, error) {
	if len(v.value) <= i {
		return nil, fmt.Errorf("index out of range: %s[%d]", v.name, i)
	}
	return v.value[i], nil
}

func (v TList) ToEnvVars(ShellFieldSeparator string) []string {
	var envVars []string
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	for _, v := range v.value {
		for _, memberEnvVar := range v.ToEnvVars(ShellFieldSeparator) {
			ev := trimmedName + "_" + memberEnvVar
			envVars = append(envVars, ev)
		}
	}
	return envVars
}

func CreateListFromInterface(name string, values []interface{}) VariableI {
	var tl TList
	tl.name = name
	tl.value = make([]VariableI, len(values))
	for i, v := range values {
		tl.value[i] = CreateVariable(strconv.Itoa(i), v)
	}
	return tl
}

func CreateListFromVars(name string, values []VariableI) VariableI {
	var tl TList
	tl.name = name
	tl.value = make([]VariableI, len(values))
	for i, v := range values {
		tl.value[i] = v
	}
	return tl
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
