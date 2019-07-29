package tvar

import (
	"strings"

	"github.com/kassybas/mate/internal/keywords"
)

type TBool struct {
	name  string
	value bool
}

func (v TBool) IsScalar() bool {
	return true
}

func (v TBool) Type() TVarType {
	return TBoolType
}

func (v TBool) Name() string {
	return v.name
}

func (v TBool) Value() interface{} {
	return v.value
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

func (v TBool) ToEnvVars() []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
