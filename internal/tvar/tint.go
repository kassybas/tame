package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/tame/types/vartype"

	"github.com/kassybas/tame/internal/keywords"
)

type TInt struct {
	TBaseVar
	value int
}

func NewInt(name string, value int) TInt {
	return TInt{
		TBaseVar: TBaseVar{
			name:     name,
			iValue:   interface{}(value),
			isScalar: true,
			varType:  vartype.TIntType,
		},
		value: value,
	}
}

func (v TInt) ToInt() (int, error) {
	return v.value, nil
}

func (v TInt) ToStr() string {
	return strconv.Itoa(v.value)
}

func (v TInt) ToEnvVars(ShellFieldSeparator string) []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
