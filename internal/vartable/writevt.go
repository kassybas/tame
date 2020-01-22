package vartable

import (
	"fmt"

	"github.com/kassybas/tame/internal/dotref"
	"github.com/kassybas/tame/internal/tvar"
)

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
		nameFields, err := dotref.ParseVarRef(names[i])
		if err != nil {
			return err
		}
		nameFields, err = vt.resolveEachField(nameFields)
		if err != nil {
			return fmt.Errorf("could not resolve variable names in fields: %s\n\t%s", names[i], err.Error())
		}
		if vt.Exists(nameFields[0].FieldName) {
			// variable exists
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
