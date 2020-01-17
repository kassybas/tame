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
	for _, p := range params {
		if vt.Exists(p.Name) {
			val, err := vt.GetVar(p.Name)
			if err != nil {
				return err
			}
			vt.Add(p.Name, val.Value())
			continue
		}
		if p.HasDefault {
			err := vt.Add(p.Name, p.DefaultValue)
			if err != nil {
				return fmt.Errorf("error adding default parameter: %s\n\t", p.Name, err.Error())
			}
			continue
		}
		return fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return nil
}
