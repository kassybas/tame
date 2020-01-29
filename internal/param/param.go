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
			// resolving default value for possible global variables
			val, err := vt.ResolveValue(p.DefaultValue)
			if err != nil {
				return fmt.Errorf("error while resolving default parameter: %s\n\t%s", p.Name, err.Error())
			}
			err = vt.Add(p.Name, val)
			if err != nil {
				return fmt.Errorf("error while creating default parameter: %s\n\t%s", p.Name, err.Error())
			}
			continue
		}
		return fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return nil
}
