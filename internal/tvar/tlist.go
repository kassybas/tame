package tvar

import (
	"fmt"
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

func CreateListFromBracketsName(name string, value interface{}) (VariableI, error) {
	lBr := strings.Index(name, keywords.IndexingSeparatorL) + 1
	rBr := strings.Index(name, keywords.IndexingSeparatorR)
	index, err := strconv.Atoi(name[lBr:rBr])
	if err != nil {
		return nil, fmt.Errorf("not integer index: %s %d", name, name[lBr:rBr])
	}
	var tl TList
	tl.name = name[0 : lBr-1]
	tl.value = make([]VariableI, index+1)
	for i := range tl.value {
		// Null all values other than the index
		tl.value[i] = TNull{}
	}
	tl.value[index] = CreateVariable(strconv.Itoa(index), value)
	fmt.Println("Created list", tl.value)
	return tl, nil
}
