package vartable

import (
	"fmt"
	"strings"

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

func (vt *VarTable) Add(name string, value interface{}) {
	v := tvar.CreateVariable(name, value)
	vt.vars[name] = v
}

func (vt *VarTable) Append(names []string, values []tvar.VariableI) {
	for i := range names {
		vt.Add(names[i], values[i])
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
	fmt.Println("---HERE IT IS")
	fmt.Println(strings.Join(formattedVars, "\n"))
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

// TODO: should get an interface in case a list or map is returned
func (vt VarTable) ResolveValue(refStr string) (tvar.VariableI, error) {
	if !strings.HasPrefix(refStr, keywords.PrefixReference) {
		// No resolution needed for constant value
		return tvar.CreateVariable("", refStr), nil
	}
	resolvedVar, err := vt.GetVar(refStr)
	return resolvedVar, err
}
