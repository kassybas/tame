package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/tame/types/vartype"

	"github.com/kassybas/tame/internal/keywords"
)

type TString struct {
	TBaseVar
	value string
}

func NewString(name string, value string) TString {
	return TString{
		TBaseVar: TBaseVar{
			name:     name,
			iValue:   interface{}(value),
			isScalar: true,
			varType:  vartype.TStringType,
		},
		value: value,
	}

}

func (v TString) ToInt() (int, error) {
	i, err := strconv.Atoi(v.value)
	return i, err
}

func (v TString) ToStr() string {
	return v.value
}

func (v TString) ToEnvVars(ShellFieldSeparator string) []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.value}
}
