package vartable

import (
	"fmt"

	"github.com/kassybas/tame/internal/texpression"
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
			continue
		}
		nameFields, err := texpression.NewExpression(names[i])
		if err != nil {
			return err
		}
		if vt.Exists(nameFields[0].FieldName) {
			// variable exists
			oldVar, err := vt.resolveFieldsVar(nameFields)
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
func (vt *VarTable) Delete(varName interface{}) error {
	s, ok := varName.(string)
	if !ok {
		return fmt.Errorf("non-string variable name: %v", varName)
	}
	if !vt.Exists(s) {
		return fmt.Errorf("variable does not exist: %v", varName)
	}
	vt.Lock()
	delete(vt.vars, s)
	vt.Unlock()
	return nil
}

func (vt *VarTable) AddVariables(newVars []tvar.TVariable) {
	for _, v := range newVars {
		vt.AddVar(v)
	}
}
