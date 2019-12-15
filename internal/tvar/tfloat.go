package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TFloat struct {
	name  string
	value float64
}

func (v TFloat) IsScalar() bool {
	return true
}

func (v TFloat) Type() vartype.TVarType {
	return vartype.TFloatType
}

func (v TFloat) Name() string {
	return v.name
}

func (v TFloat) Value() interface{} {
	return v.value
}

func (v TFloat) ToInt() (int, error) {
	return int(v.value), nil
}

func (v TFloat) ToStr() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

func (v TFloat) ToEnvVars() []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
