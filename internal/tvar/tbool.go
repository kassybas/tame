package tvar

import (
	"strings"

	"github.com/kassybas/tame/types/vartype"

	"github.com/kassybas/tame/internal/keywords"
)

type TBool struct {
	TBaseVar
	value bool
}

func NewBool(name string, value bool) TBool {
	return TBool{
		TBaseVar: TBaseVar{
			name:     name,
			iValue:   interface{}(value),
			isScalar: true,
			varType:  vartype.TBoolType,
		},
		value: value,
	}
}

func (v TBool) ToInt() (int, error) {
	if v.value {
		return 1, nil
	}
	return 0, nil
}

func (v TBool) ToStr() string {
	if v.value {
		return "true"
	}
	return "false"
}

func (v TBool) ToEnvVars(ShellFieldSeparator string) []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
