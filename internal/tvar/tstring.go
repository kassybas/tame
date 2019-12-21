package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TString struct {
	name  string
	value string
}

func (v TString) IsScalar() bool {
	return true
}

func (v TString) Type() vartype.TVarType {
	return vartype.TStringType
}

func (v TString) Name() string {
	return v.name
}

func (v TString) Value() interface{} {
	return v.value
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
