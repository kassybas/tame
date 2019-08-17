package vartable

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kassybas/mate/internal/dotref"
	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/tvar"
)

type VarTable struct {
	vars map[string]tvar.VariableI
}

func (vt VarTable) GetVar(name string) (tvar.VariableI, error) {
	val, exists := vt.vars[name]
	if !exists {
		return nil, fmt.Errorf("variable does not exist:'%s'", name)
	}
	return val, nil
}

func (vt VarTable) Exists(name string) bool {
	_, exists := vt.vars[name]
	return exists
}

func NewVarTable() VarTable {
	vt := VarTable{}
	vt.vars = make(map[string]tvar.VariableI)
	return vt
}

func (vt *VarTable) Add(v tvar.VariableI) {
	oldVar, err := vt.GetVar(v.Name())
	// if already exists
	if err == nil {
		if oldVar.Type() == tvar.TMapType {
			v = tvar.UpdateCompositeValue(oldVar, v)
		}
		if oldVar.Type() == tvar.TListType && v.Type() == tvar.TListType {
			v, err = tvar.MergeLists(oldVar.(tvar.TList), v.(tvar.TList))
			if err != nil {
				logrus.Fatal(err)
			}
		}
	}
	vt.vars[v.Name()] = v
}

func (vt *VarTable) Append(names []string, values []interface{}) {
	for i := range names {
		v := tvar.CreateVariable(names[i], values[i])
		vt.Add(v)
	}
}

func (vt *VarTable) AddVariables(newVars []tvar.VariableI) {
	for _, v := range newVars {
		vt.vars[v.Name()] = v
	}
}

func (vt *VarTable) GetAllEnvVars() []string {
	formattedVars := []string{}
	for _, v := range vt.vars {
		if v != nil {
			formattedVars = append(formattedVars, v.ToEnvVars()...)
		}
	}
	return formattedVars
}

func (vt VarTable) ResolveVar(v tvar.VariableI) (tvar.VariableI, error) {
	// TODO: handle composite types
	value, err := vt.ResolveValue(v.ToStr())
	if err != nil {
		return nil, err
	}
	return tvar.CreateVariable(v.Name(), value), nil
}

func (vt VarTable) GetVarByDotRef(dr dotref.DotRef) (tvar.VariableI, error) {
	cur, err := vt.GetVar(dr.Name)
	if err != nil {
		return nil, err
	}

	for _, field := range dr.Fields {
		if field.FieldName == "" {
			// list reference
			if cur.Type() != tvar.TListType {
				return nil, fmt.Errorf("indexing non-list type: %d -- %s", cur.Name, field.Index)
			}
			cur, err = cur.(tvar.TList).GetItem(field.Index)
			if err != nil {
				return nil, err
			}
			continue
		}
		// map reference
		if cur.Type() != tvar.TMapType {
			return nil, fmt.Errorf("field reference on a non-map type: %s -- %s", cur.Name, field.FieldName)
		}
		cur, err = cur.(tvar.TMap).GetMember(field.FieldName)
		if err != nil {
			return nil, err
		}
	}
	return cur, nil
}

func (vt VarTable) ResolveValue(refStr string) (interface{}, error) {
	if !strings.HasPrefix(refStr, keywords.PrefixReference) {
		// No resolution needed for constant value
		return refStr, nil
	}
	dr, err := dotref.NewReference(refStr, nil)
	if err != nil {
		return nil, err
	}
	value, err := vt.GetVarByDotRef(dr)
	if err != nil {
		return nil, fmt.Errorf("%s\n\tfailed to resolve variable reference: %s", err, refStr)

	}
	return value.Value(), err
}
