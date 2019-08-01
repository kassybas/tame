package tvar

import (
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

func (tm TMap) UpdateMap(newMap TMap) TMap {
	for k, v := range newMap.value {
		if !tm.IsMember(k) {
			tm.value[k] = v
			return tm
		}
		if tm.value[k].Type() != TMapType {
			tm.value[k] = v
			return tm
		}
		tm.value[k].(TMap).UpdateMap(v.(TMap))
	}
	return tm
}
