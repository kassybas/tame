package target

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/steprunner"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/opts"
)

type Param struct {
	Name         string
	HasDefault   bool
	DefaultValue interface{}
}

type Target struct {
	Ctx       *tcontext.Context
	Name      string
	Steps     []step.Step
	Params    []Param
	Opts      opts.ExecutionOpts
	Variables []tvar.TVariable
	Summary   string
	Status    int
}

func (t Target) Make(vt *vartable.VarTable, parentOpts opts.ExecutionOpts) step.StepStatus {
	vt.AddVariables(t.Ctx.Globals)
	err := resolveParams(vt, t.Params)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("could not resolve parameters in target: %s\n\t%s", t.Name, err)}
	}
	// inherit silent
	t.Opts.Silent = parentOpts.Silent

	status := steprunner.RunAllSteps(t.Steps, *t.Ctx, vt, t.Opts)
	return status
}

func resolveParams(vt *vartable.VarTable, params []Param) error {
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
			vt.Add(p.Name, p.DefaultValue)
			continue
		}
		return fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return nil
}

func (t Target) IsParameter(name string) bool {
	for _, p := range t.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
