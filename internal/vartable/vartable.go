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
		if names[i] == "" {
			// ignored
			continue
		}
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

func (vt *VarTable) GetAllValues() map[string]interface{} {
	vars := make(map[string]interface{})
	for k, v := range vt.vars {
		vars[k] = v.Value()
	}
	return vars
}

func (vt VarTable) ResolveValue(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case string:
		{
			if !strings.HasPrefix(val, keywords.PrefixReference) {
				// No resolution needed for constant value
				return val, nil
			}
			valueVar, err := vt.GetVar(val)
			if err != nil {
				return nil, err
			}
			if !dotref.IsDotRef(val) {
				return valueVar.Value(), nil
			}
			fields, err := dotref.ParseFields(val)
			if err != nil {
				return nil, err
			}
			// resolve each field
			fields, err = vt.resolveEachField(fields)
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

func (vt VarTable) resolveEachField(fields []dotref.RefField) ([]dotref.RefField, error) {
	if len(fields) < 2 {
		// single field does not need resolve
		return fields, nil
	}
	for i := range fields {
		if i == 0 {
			// field 0 is the variable name, it would resolve to itself
			continue
		}
		if !strings.HasPrefix(fields[i].FieldName, keywords.PrefixReference) {
			// No resolution needed for constant value
			continue
		}

		resolvedField, err := vt.ResolveValue(fields[i].FieldName)
		if err != nil {
			return nil, err
		}
		switch resolvedField := resolvedField.(type) {
		case int:
			{
				fields[i].FieldName = ""
				fields[i].Index = resolvedField
			}
		case string:
			{
				fields[i].FieldName = resolvedField
			}
		default:
			{
				return nil, fmt.Errorf("unknown field reference or index variable type: can only be string or int, got: %v :: %T", resolvedField, resolvedField)
			}
		}
	}
	return fields, nil
}
