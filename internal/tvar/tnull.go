package tvar

import (
	"github.com/kassybas/tame/types/vartype"
)

type TNull struct {
	TBaseVar
}

func NewNull(name string) TNull {
	return TNull{
		TBaseVar: TBaseVar{
			iValue:   interface{}(nil),
			name:     name,
			isScalar: false,
			varType:  vartype.TNullType,
		},
	}
}

func (v TNull) ToInt() (int, error) {
	return 0, nil
}

func (v TNull) ToStr() string {
	return ""
}

func (v TNull) ToEnvVars(ShellFieldSeparator string) []string {
	return []string{}
}
