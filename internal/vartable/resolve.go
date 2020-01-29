package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/texpression"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/exprtype"
)

func (vt *VarTable) resolveFieldsVar(refFields []texpression.ExprField) (tvar.TVariable, error) {
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
	return vt.getVarByFields(refFields)
}

func (vt *VarTable) resolveFieldsValue(refFields []texpression.ExprField) (interface{}, error) {
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

func (vt *VarTable) resolveValueString(val string) (interface{}, error) {
	if !strings.HasPrefix(val, keywords.PrefixReference) {
		// No resolution needed for constant value
		if strings.HasPrefix(val, "\\$") {
			// remove escape sign from before $
			val = strings.TrimPrefix(val, "\\")
		}
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

func (vt *VarTable) resolveValueMap(val map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	for k, v := range val {
		resK, err := vt.ResolveValue(k)
		if err != nil {
			return nil, err
		}
		if resK != k {
			vt.Delete(k)
		}
		val[resK], err = vt.ResolveValue(v)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

func (vt *VarTable) resolveValueList(val []interface{}) ([]interface{}, error) {
	var err error
	for i := range val {
		val[i], err = vt.ResolveValue(val[i])
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

func (vt *VarTable) ResolveValue(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case string:
		{
			return vt.resolveValueString(val)
		}
	case map[interface{}]interface{}:
		{
			return vt.resolveValueMap(val)
		}
	case []interface{}:
		{
			return vt.resolveValueList(val)
		}
	default:
		{
			return val, nil
		}
	}
}
