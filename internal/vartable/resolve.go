package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/eval"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/texpression"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/exprtype"
)

func (vt *VarTable) EvaluateExpression(expression string) (interface{}, error) {
	env := vt.GetAllValues()
	return eval.EvaluateExpression(expression, env)
}

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
		if refFields[i].Type == exprtype.Expression {
			res, err := vt.EvaluateExpression(refFields[i].Val)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve expression: %s\n\t%s", refFields[i].Val, err.Error())
			}
			refFields[i], err = texpression.NewField(res)
			if err != nil {
				return nil, err
			}
		}
	}
	return vt.getVarByFields(refFields)
}

func (vt *VarTable) resolveSingleField(field texpression.ExprField) (interface{}, error) {
	switch field.Type {
	case exprtype.Expression:
		res, err := vt.EvaluateExpression(field.Val)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve expression: %s\n\t%s", field.Val, err.Error())
		}
		return res, nil
	case exprtype.Index:
		return field.Index, nil
	case exprtype.Literal:
		return field.Val, nil
	default:
		return nil, fmt.Errorf("internal error: could not resolve single field: %+v", field)
	}
}

func (vt *VarTable) resolveFieldsValue(refFields []texpression.ExprField) (interface{}, error) {
	if len(refFields) == 1 && refFields[0].Type != exprtype.VarName {
		// last field is end of recursion
		return vt.resolveSingleField(refFields[0])
	}
	// resolveFields resolves the fields in refFields in place
	// note: refFields slice is passed by reference
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
	// requires resolution
	if strings.HasPrefix(val, keywords.PrefixReference) || strings.HasPrefix(val, "(") {
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
	// No resolution needed for constant value
	// In literals remove escape sign from before $ and (
	if strings.HasPrefix(val, "\\$") || strings.HasPrefix(val, "\\(") {
		val = strings.TrimPrefix(val, "\\")
	}
	return val, nil
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

func (vt *VarTable) ResolveValueToStr(val interface{}) (string, error) {
	v, err := vt.ResolveValue(val)
	if err != nil {
		return "", err
	}
	sVal, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("expression does not resolve to string: %v", val)
	}
	return sVal, nil
}
