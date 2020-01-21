package vartable

import (
	"fmt"
	"strings"
	"sync"

	"github.com/antonmedv/expr"
	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/tvar"
)

type VarTable struct {
	sync.RWMutex
	vars map[string]tvar.TVariable
}

func CopyVarTable(vt *VarTable) *VarTable {
	newVt := VarTable{
		vars: make(map[string]tvar.TVariable),
	}
	vt.RLock()
	for k, v := range vt.vars {
		newVt.vars[k] = v
	}
	vt.RUnlock()
	return &newVt
}

func NewVarTable() *VarTable {
	vt := VarTable{}
	vt.vars = make(map[string]tvar.TVariable)
	return &vt
}

func (vt *VarTable) GetVar(fullName string) (tvar.TVariable, error) {
	fields, err := dotref.ParseFields(fullName)
	if err != nil || fields[0].FieldName == "" {
		return nil, fmt.Errorf("failed to parse variable name:%s\n%s", fullName, err)
	}
	return vt.GetVarByFields(fields)
}

func (vt *VarTable) GetVarByFields(fields []dotref.RefField) (tvar.TVariable, error) {
	name := fields[0].FieldName
	vt.RLock()
	val, exists := vt.vars[name]
	vt.RUnlock()
	if !exists {
		return nil, fmt.Errorf("variable '%s' does not exist", name)
	}
	return val, nil
}

func (vt *VarTable) Exists(name string) bool {
	vt.RLock()
	_, exists := vt.vars[name]
	vt.RUnlock()
	return exists
}

func (vt *VarTable) AddVar(v tvar.TVariable) error {
	if v == nil {
		return nil
	}
	vt.Lock()
	vt.vars[v.Name()] = v
	vt.Unlock()
	return nil
}

func (vt *VarTable) Add(name string, value interface{}) error {
	if name == "" {
		return nil
	}
	newVar := tvar.NewVariable(name, value)
	vt.Lock()
	vt.vars[name] = newVar
	vt.Unlock()
	return nil
}

func (vt *VarTable) Append(names []string, values []interface{}) error {
	for i := range names {
		if names[i] == "" {
			// ignored
			continue
		}
		nameFields, err := dotref.ParseFields(names[i])
		if err != nil {
			return err
		}
		nameFields, err = vt.resolveEachField(nameFields)
		if err != nil {
			return fmt.Errorf("could not resolve variable names in fields: %s\n\t%s", names[i], err.Error())
		}
		if vt.Exists(nameFields[0].FieldName) {
			// variable exists
			//
			oldVar, err := vt.GetVarByFields(nameFields)
			if err != nil {
				return err
			}
			newVar, err := oldVar.SetValue(nameFields, values[i])
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
		vt.AddVar(v)
	}
}

func (vt *VarTable) GetAllEnvVars(ShellFieldSeparator string) []string {
	formattedVars := []string{}
	vt.RLock()
	allVars := vt.vars
	vt.RUnlock()
	for _, v := range allVars {
		formattedVars = append(formattedVars, v.ToEnvVars(ShellFieldSeparator)...)
	}
	return formattedVars
}

func (vt *VarTable) GetAllValues() map[string]interface{} {
	vars := make(map[string]interface{})
	vt.RLock()
	allVars := vt.vars
	vt.RUnlock()
	for k, v := range allVars {
		vars[k] = v.Value()
	}
	return vars
}

func (vt *VarTable) ResolveValue(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case string:
		{
			if !strings.HasPrefix(val, keywords.PrefixReference) {
				// No resolution needed for constant value
				return val, nil
			}
			env := vt.GetAllValues()
			program, err := expr.Compile(val, expr.Env(env))
			if err != nil {
				return nil, fmt.Errorf("could not parse variable reference: %s", err.Error())
			}
			result, err := expr.Run(program, env)
			if err != nil {
				return nil, fmt.Errorf("could not resolve variable reference: %s", err.Error())
			}
			return result, nil
			// valueVar, err := vt.GetVar(val)
			// if err != nil {
			// 	return nil, err
			// }
			// if !dotref.IsDotRef(val) {
			// 	return valueVar.Value(), nil
			// }
			// fields, err := dotref.ParseFields(val)
			// if err != nil {
			// 	return nil, err
			// }
			// // resolve each field
			// fields, err = vt.resolveEachField(fields)
			// if err != nil {
			// 	return nil, err
			// }
			// return valueVar.GetInnerValue(fields)
		}
	default:
		{
			return val, nil
		}
	}
}

func (vt *VarTable) resolveEachField(fields []dotref.RefField) ([]dotref.RefField, error) {
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
