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
			ev := trimmedName + "_" + memberEnvVar
			envVars = append(envVars, ev)
			fmt.Println("Added env var:", ev)
		}
	}
	return envVars
}

func CreateMap(name string, value interface{}) VariableI {
	var tm TMap
	tm.name = name
	tm.value = make(map[string]VariableI)
	switch value.(type) {
	case map[interface{}]interface{}:
		{
			for k, v := range value.(map[interface{}]interface{}) {
				tm.value[k.(string)] = CreateVariable(k.(string), v)
			}
		}
	}
	return tm
}
