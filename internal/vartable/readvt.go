package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/dotref/reftype"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/tvar"
)

func (vt *VarTable) GetVar(fullName string) (tvar.TVariable, error) {
	fields, err := dotref.ParseVarRef(fullName)
	if err != nil || fields[0].FieldName == "" {
		return nil, fmt.Errorf("failed to parse variable name:%s\n%s", fullName, err)
	}
	return vt.GetVarByFields(fields)
}

func (vt *VarTable) GetValueByFields(fields []dotref.RefField) (interface{}, error) {
	v, err := vt.GetVarByFields(fields)
	if err != nil {
		return nil, err
	}
	val, err := v.GetInnerValue(fields)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (vt *VarTable) GetVarByFields(fields []dotref.RefField) (tvar.TVariable, error) {
	name := fields[0].FieldName
	strName, ok := name.(string)
	if !ok {
		return nil, fmt.Errorf("non-string variable name: %s (type %T)", name, name)
	}
	vt.RLock()
	val, exists := vt.vars[strName]
	vt.RUnlock()
	if !exists {
		return nil, fmt.Errorf("variable '%s' does not exist", name)
	}
	return val, nil
}

func (vt *VarTable) ResolveToField(fields []dotref.RefField) (dotref.RefField, error) {
	val, err := vt.GetValueByFields(fields)
	if err != nil {
		return dotref.RefField{}, err
	}

	switch val := val.(type) {
	case int:
		return dotref.RefField{
			Index: val,
			Type:  reftype.Index,
		}, nil
	default:
		return dotref.RefField{
			FieldName: val,
			Type:      reftype.Literal,
		}, nil
	}
	// TODO: innerfields?

}

func (vt *VarTable) ResolveRefValue(fields []dotref.RefField) (interface{}, error) {
	valueVar, err := vt.GetVarByFields(fields)
	if err != nil {
		return nil, err
	}
	// resolve each field
	// fields, err = vt.resolveEachField(fields)
	for i := range fields {
		if fields[i].Type == reftype.InnerRef {
			fields[i], err = vt.ResolveToField(fields[i].InnerRefs)
			if err != nil {
				return nil, err
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return valueVar.GetInnerValue(fields)
}

func (vt *VarTable) ResolveValue(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case string:
		{
			if !strings.HasPrefix(val, keywords.PrefixReference) {
				// No resolution needed for constant value
				return val, nil
			}
			fields, err := dotref.ParseVarRef(val)
			if err != nil {
				return nil, err
			}
			return vt.ResolveRefValue(fields)

		}
	default:
		{
			return val, nil
		}
	}
}

// func (vt )

// func (vt *VarTable) resolveEachField(fields []dotref.RefField) ([]dotref.RefField, error) {
// 	var err error
// 	for i := range fields {
// 		if i == 0 {
// 			// we are looking for the field of the map, to find it we resolve
// 			// each variable field in the name
// 			// field 0 is the variable name, it would resolve to itself
// 			continue
// 		}
// 		switch fields[i].Type {
// 		case reftype.Literal:
// 			// No resolution needed for constant value
// 			continue
// 		case reftype.VarName:
// 			return nil, fmt.Errorf("illegal variable in dot reference")
// 		case reftype.InnerRef:
// 			resolvedInnerValue, err := vt.resolveEachField(fields[i].InnerRefs)
// 			if err != nil {
// 				return nil, fmt.Errorf("could not resolve variable: %v\n\t%s", field[i], err.Error())
// 			}
// 			fields[i].Type = reftype.Literal
// 			fields[i].InnerRefs = nil
// 		}

// 		// resolvedField, err := vt.ResolveValue(fields[i].FieldName)
// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 		// switch resolvedField := resolvedField.(type) {
// 		// case int:
// 		// 	{
// 		// 		fields[i].FieldName = ""
// 		// 		fields[i].Index = resolvedField
// 		// 	}
// 		// case string:
// 		// 	{
// 		// 		fields[i].FieldName = resolvedField
// 		// 	}
// 		// default:
// 		// 	{
// 		// 		return nil, fmt.Errorf("unknown field reference or index variable type: can only be string or int, got: %v :: %T", resolvedField, resolvedField)
// 		// 	}
// 		// }
// 	}
// 	return fields, nil
// }
