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
		fields, err := dotref.ParseFields(names[i])
		if err != nil {
			return err
		}
		if oldVar != nil {
			// variable exists
			newVar, err := oldVar.SetValue(fields, values[i])
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

func (vt VarTable) ResolveValue(val interface{}) (interface{}, error) {
	switch val.(type) {
	case string:
		{
			s := val.(string)
			if !strings.HasPrefix(s, keywords.PrefixReference) {
				// No resolution needed for constant value
				return val, nil
			}
			valueVar, err := vt.GetVar(s)
			if err != nil {
				return nil, err
			}
			if !dotref.IsDotRef(s) {
				return valueVar.Value(), nil
			}
			fields, err := dotref.ParseFields(s)
			if err != nil {
				return nil, err
			}
			return valueVar.GetInnerValue(fields)
		}
	default:
		{
			return val, nil
		}
	}
}
