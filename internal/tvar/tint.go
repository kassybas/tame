package tvar

import (
	"strconv"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
)

type TInt struct {
	name  string
	value int
}

func (v TInt) IsScalar() bool {
	return true
}

func (v TInt) Type() TVarType {
	return TIntType
}

func (v TInt) Name() string {
	return v.name
}

func (v TInt) Value() interface{} {
	return v.value
}

func (v TInt) ToInt() (int, error) {
	return v.value, nil
}

func (v TInt) ToStr() string {
	return strconv.Itoa(v.value)
}

func (v TInt) ToEnvVars() []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
