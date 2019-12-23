package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/tvar"
)

type VarTable struct {
	vars map[string]tvar.TVariable
}

func (vt VarTable) GetVar(fullName string) (tvar.TVariable, error) {
	fields, err := dotref.ParseFields(fullName)
	if err != nil || fields[0].FieldName == "" {
		return nil, fmt.Errorf("failed to parse variable name:%s\n%s", fullName, err)
	}
	name := fields[0].FieldName
	val, exists := vt.vars[name]
	if !exists {
		return nil, nil
	}
	return val, nil
}

func (vt VarTable) Exists(name string) bool {
	_, exists := vt.vars[name]
	return exists
}

func NewVarTable() VarTable {
	vt := VarTable{}
	vt.vars = make(map[string]tvar.TVariable)
	return vt
}

func (vt *VarTable) AddVar(v tvar.TVariable) error {
	vt.vars[v.Name()] = v
	return nil
}

func (vt *VarTable) Add(name string, value interface{}) error {
	vt.vars[name] = tvar.NewVariable(name, value)
	return nil
}

func (vt *VarTable) Append(names []string, values []interface{}) error {
	for i := range names {
		oldVar, err := vt.GetVar(names[i])
		if err != nil {
			return err
		}
		dr, err := dotref.ParseFields(names[i])
		if err != nil {
			return err
		}
		if oldVar != nil {
			// variable exists
			if len(dr) > 1 {
				return fmt.Errorf("assingment to undeclared map: %s", names[i])
			}
			newVar, err := oldVar.SetValue(dr, values[i])
			if err != nil {
				return err
			}
			vt.AddVar(newVar)
		} else {
			// new variable
			vt.Add(names[i], values[i])
		}
	}
	return nil
}

func (vt *VarTable) AddVariables(newVars []tvar.TVariable) {
	for _, v := range newVars {
		vt.vars[v.Name()] = v
	}
}

func (vt *VarTable) GetAllEnvVars(ShellFieldSeparator string) []string {
	formattedVars := []string{}
	for _, v := range vt.vars {
		formattedVars = append(formattedVars, v.ToEnvVars(ShellFieldSeparator)...)
	}
	return formattedVars
}

func (vt VarTable) ResolveVar(v tvar.TVariable) (tvar.TVariable, error) {
	if !strings.HasPrefix(v.ToStr(), keywords.PrefixReference) {
		return v, nil
	}
	value, err := vt.ResolveValue(v.ToStr())
	if err != nil {
		return nil, err
	}
	v.SetValue([]dotref.RefField{}, value)
	return v, nil
}

// func (vt VarTable) GetVarByDotRef(dr dotref.DotRef) (tvar.TVariable, error) {
// 	cur, err := vt.GetVar(dr.Name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, field := range dr.Fields {
// 		if field.FieldName == "" {
// 			// list reference
// 			if cur.Type() != vartype.TListType {
// 				return nil, fmt.Errorf("indexing non-list type: %s[%d] (type: %s)", cur.Name(), field.Index, vartype.GetTypeNameString(cur.Type()))
// 			}
// 			cur, err = cur.(tvar.TList).GetItem(field.Index)
// 			if err != nil {
// 				return nil, err
// 			}
// 			continue
// 		}
// 		// map reference
// 		if cur.Type() != vartype.TMapType {
// 			return nil, fmt.Errorf("field reference on a non-map type: %s.%s (type: %s)", cur.Name(), field.FieldName, vartype.GetTypeNameString(cur.Type()))
// 		}
// 		cur, err = cur.(tvar.TMap).GetMember(field.FieldName)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return cur, nil
// }

func (vt VarTable) ResolveValue(refStr string) (interface{}, error) {
	if !strings.HasPrefix(refStr, keywords.PrefixReference) {
		// No resolution needed for constant value
		return refStr, nil
	}
	value, err := vt.GetVar(refStr)
	if err != nil {
		return nil, fmt.Errorf("%s\n\tfailed to resolve variable reference: %s", err, refStr)
	}
	return value.Value(), err
}
