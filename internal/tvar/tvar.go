package tvar

import (
	"strings"

	"github.com/kassybas/mate/internal/keywords"
)

type Variable struct {
	Name string
	// TODO: interface
	Value string
}

func CreateVariable(name string, value interface{}) Variable {
	return Variable{
		Name:  name,
		Value: value.(string),
	}
}

func (v Variable) FormatToEnvVars() []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.Name, keywords.PrefixReference)
	// TODO: handle list + maps
	return []string{trimmedName + "=" + v.Value}
}
