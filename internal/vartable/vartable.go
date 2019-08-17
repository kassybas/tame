package vartable

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

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

func (vt *VarTable) AddVar(newVar tvar.VariableI) {
	vt.vars[newVar.Name()] = newVar
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
	if !strings.HasPrefix(v.ToStr(), keywords.PrefixReference) {
		// No resolution needed for constant value
		return v, nil
	}
	resolvedVar, err := vt.GetVar(v.ToStr())
	if err != nil {
		return nil, err
	}
	return tvar.CreateVariable(v.Name(), resolvedVar.Value()), nil
}

func (vt VarTable) ResolveValue(refStr string) (interface{}, error) {
	if !strings.HasPrefix(refStr, keywords.PrefixReference) {
		// No resolution needed for constant value
		return tvar.CreateVariable("", refStr).Value(), nil
	}
	resolvedVar, err := vt.GetVar(refStr)
	return resolvedVar.Value(), err
}
