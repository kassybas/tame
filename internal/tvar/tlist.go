package tvar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/types/vartype"
)

type TList struct {
	name  string
	value []VariableI
}

func (v TList) Type() vartype.TVarType {
	return vartype.TListType
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

func (v TList) GetItem(i int) (VariableI, error) {
	if len(v.value) <= i {
		return nil, fmt.Errorf("index out of range: %s[%d]", v.name, i)
	}
	return v.value[i], nil
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
