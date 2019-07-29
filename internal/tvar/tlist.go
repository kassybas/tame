package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
)

type TList struct {
	name  string
	value []VariableI
}

func (v TList) Type() TVarType {
	return TListType
}

func (v TList) IsScalar() bool {
	return false
}

func (v TList) Name() string {
	return v.name
}

func (v TList) Value() interface{} {
	return v.value
}

func (v TList) ToInt() (int, error) {
	// i, err := strconv.Atoi(v.value.(string))
	// return i, err
	return 0, nil
}

func (v TList) ToStr() string {
	return ""
	// return v.value.(string)
}

func (v TList) ToEnvVars() []string {
	var envVars []string
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	for _, v := range v.value {
		for _, memberEnvVar := range v.ToEnvVars() {
			ev := trimmedName + "_" + memberEnvVar
			envVars = append(envVars, ev)
		}
	}
	return envVars
}

func CreateList(name string, value []interface{}) VariableI {
	var tl TList
	tl.name = name
	tl.value = make([]VariableI, len(value))
	for i, elem := range value {
		tl.value[i] = CreateVariable(strconv.Itoa(i), elem)
	}
	return tl
}
