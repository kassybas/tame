package vartable

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/texpression"
	"github.com/kassybas/tame/internal/tvar"
)

func (vt *VarTable) GetVar(fullName string) (tvar.TVariable, error) {
	fields, err := texpression.NewExpression(fullName)
	if err != nil || fields[0].Val == "" {
		return nil, fmt.Errorf("failed to parse variable name:%s\n%s", fullName, err)
	}
	return vt.getVarByFields(fields)
}

func (vt *VarTable) getVarByName(name string) (tvar.TVariable, error) {
	vt.RLock()
	val, exists := vt.vars[name]
	vt.RUnlock()
	if !exists {
		msg := fmt.Sprintf("variable '%s' does not exist", name)
		if !strings.HasPrefix(name, keywords.PrefixReference) {
			msg = fmt.Sprintf("%s: variable should start with $ '$%s', or quoted if literal in field '\"%s\"'", msg, name, name)
		}
		return nil, fmt.Errorf(msg)
	}
	return val, nil
}

func (vt *VarTable) getVarByFields(fields []texpression.ExprField) (tvar.TVariable, error) {
	return vt.getVarByName(fields[0].Val)
}
