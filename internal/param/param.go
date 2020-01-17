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

func ResolveParams(vt *vartable.VarTable, params []Param) error {
	newVt := vartable.NewVarTable()
	for _, p := range params {
		if vt.Exists(p.Name) {
			val, err := vt.GetVar(p.Name)
			if err != nil {
				return err
			}
			newVt.Add(p.Name, val.Value())
			continue
		}
		if p.HasDefault {
			newVt.Add(p.Name, p.DefaultValue)
			continue
		}
		return fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return nil
}
