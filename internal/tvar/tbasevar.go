package tvar

import (
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
