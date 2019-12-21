package tvar

import (
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/vartype"
)

type TBaseVar struct {
	name     string
	iValue   interface{}
	isScalar bool
	varType  vartype.TVarType
}

func (v TBaseVar) Name() string {
	return v.name
}

func (v TBaseVar) Value() interface{} {
	return v.iValue
}

func (v TBaseVar) IsScalar() bool {
	return v.isScalar
}

func (v TBaseVar) Type() vartype.TVarType {
	return v.varType
}

func (v TBaseVar) ToEnvVars(ShellFieldSeparator string) []string {
	// Remove $ for shell env format
	trimmedName := strings.TrimPrefix(v.name, keywords.PrefixReference)
	return []string{trimmedName + "=" + v.ToStr()}
}
