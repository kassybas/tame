package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
)

type TString struct {
	name  string
	value string
}

func (v TString) Type() TVarType {
	return TStringType
}

func (v TString) Name() string {
	return v.name
}

func (v TString) Value() interface{} {
	return v.value
}

func (v TString) ToInt() int {
	i, err := strconv.Atoi(v.value)
	if err != nil {
		return 0
	}
	return i
}

func (v TString) ToStr() string {
	return v.value
}

func (v TString) ToEnvVars() []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.value}
}
