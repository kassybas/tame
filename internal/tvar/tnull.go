package tvar

import "github.com/kassybas/mate/types/vartype"

type TNull struct {
	name string
}

func (v TNull) Type() vartype.TVarType {
	return vartype.TNullType
}

func (v TNull) IsScalar() bool {
	return false
}

func (v TNull) Name() string {
	return v.name
}

func (v TNull) Value() interface{} {
	return nil
}

func (v TNull) ToInt() (int, error) {
	return 0, nil
}

func (v TNull) ToStr() string {
	return ""
}

func (v TNull) ToEnvVars() []string {
	return []string{}
}
