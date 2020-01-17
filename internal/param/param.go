package param

import (
	"fmt"

	"github.com/kassybas/tame/internal/vartable"
)

type Param struct {
	Name         string
	HasDefault   bool
	DefaultValue interface{}
}

func ResolveParams(vt *vartable.VarTable, params []Param) (*vartable.VarTable, error) {
	newVt := vartable.NewVarTable()
	for _, p := range params {
		if vt.Exists(p.Name) {
			val, err := vt.GetVar(p.Name)
			if err != nil {
				return nil, err
			}
			newVt.Add(p.Name, val.Value())
			continue
		}
		if p.HasDefault {
			newVt.Add(p.Name, p.DefaultValue)
			continue
		}
		return nil, fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return newVt, nil
}
