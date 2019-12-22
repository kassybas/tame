package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/tame/types/vartype"

	"github.com/kassybas/tame/internal/keywords"
)

type TFloat struct {
	TBaseVar
	value float64
}

func NewFloat(name string, value float64) TFloat {
	return TFloat{
		TBaseVar: TBaseVar{
			name:     name,
			iValue:   interface{}(value),
			isScalar: true,
			varType:  vartype.TFloatType,
		},
		value: value,
	}
}

func (v TFloat) ToInt() (int, error) {
	return int(v.value), nil
}

func (v TFloat) ToStr() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

func (v TFloat) ToEnvVars(ShellFieldSeparator string) []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
