package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/texpression"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/exprtype"
)

func (vt *VarTable) resolveFieldsVar(refFields []texpression.RefField) (tvar.TVariable, error) {
	for i := range refFields {
		if refFields[i].Type == exprtype.InnerRef {
			innerVal, err := vt.resolveFieldsValue(refFields[i].InnerRefs)
			if err != nil {
				return nil, err
			}
			refFields[i], err = texpression.NewField(innerVal)
			if err != nil {
				return nil, err
			}
		}
	}
	return vt.GetVarByFields(refFields)
}

func (vt *VarTable) resolveFieldsValue(refFields []texpression.RefField) (interface{}, error) {
	v, err := vt.resolveFieldsVar(refFields)
	if err != nil {
		return nil, err
	}
	val, err := v.GetInnerValue(refFields)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (vt *VarTable) ResolveValue(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case string:
		{
			if !strings.HasPrefix(val, keywords.PrefixReference) {
				// No resolution needed for constant value
				return val, nil
			}
			fields, err := texpression.NewExpression(val)
			if err != nil {
				return nil, err
			}
			res, err := vt.resolveFieldsValue(fields)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve value: %v\n\t%s", val, err.Error())
			}
			return res, nil
		}
	default:
		{
			return val, nil
		}
	}
}
