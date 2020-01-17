package target

import (
	"github.com/kassybas/tame/internal/param"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/stepblock"
	"github.com/kassybas/tame/internal/steprunner"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/opts"
)

type Target struct {
	Ctx       *tcontext.Context
	Name      string
	Steps     stepblock.StepBlock
	Params    []param.Param
	Opts      opts.ExecutionOpts
	Variables []tvar.TVariable
	Summary   string
	Status    int
}

func (t Target) Make(vt *vartable.VarTable, parentOpts opts.ExecutionOpts) step.StepStatus {
	vt.AddVariables(t.Ctx.Globals)
	// inherit silent
	t.Opts.Silent = parentOpts.Silent

	status := steprunner.RunAllSteps(t.Steps, *t.Ctx, vt, t.Opts)
	return status
}

func (t Target) IsParameter(name string) bool {
	for _, p := range t.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
