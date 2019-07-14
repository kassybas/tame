package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/tvar"
)

type VarTable struct {
	vars map[string]tvar.Variable
}

func (vt VarTable) GetValue(name string) (tvar.Variable, error) {
	val, exists := vt.vars[name]
	if !exists {
		return tvar.Variable{}, fmt.Errorf("variable does not exist:'%s'", name)
	}
	return val, nil
}

func (vt VarTable) Exists(name string) bool {
	_, exists := vt.vars[name]
	return exists
}

func NewVarTable() VarTable {
	vt := VarTable{}
	vt.vars = make(map[string]tvar.Variable)
	return vt
}

func (vt *VarTable) Add(name string, value interface{}) {
	vt.vars[name] = tvar.CreateVariable(name, value)
}

func (vt *VarTable) AddVar(newVar tvar.Variable) {
	vt.vars[newVar.Name] = newVar
}

func (vt *VarTable) AddVariables(newVars []tvar.Variable) {
	for _, v := range newVars {
		vt.vars[v.Name] = v
	}
}

func (vt *VarTable) GetAllEnvVars() []string {
	formattedVars := []string{}
	for _, v := range vt.vars {
		// Remove $ for shell env format
		formattedVars = append(formattedVars, v.FormatToEnvVars()...)
	}
	return formattedVars

}

func (vt VarTable) ResolveVar(v tvar.Variable) (tvar.Variable, error) {
	if !strings.HasPrefix(v.Value, keywords.PrefixReference) {
		// No resolution needed for constant value
		return v, nil
	}

	resolvedVar, err := vt.GetValue(v.Value)
	if err != nil {
		return tvar.Variable{}, err
	}
	return tvar.Variable{Name: v.Name, Value: resolvedVar.Value}, nil
}
func (vt VarTable) ResolveValue(val string) (string, error) {
	if !strings.HasPrefix(val, keywords.PrefixReference) {
		// No resolution needed for constant value
		return val, nil
	}

	resolvedVar, err := vt.GetValue(val)
	return resolvedVar.Value, err
}
