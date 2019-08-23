package step

import (
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/internal/vartable"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
)

type ReturnStep struct {
	Name      string
	Arguments []tvar.VariableI
	Return    []string
}

func (s ReturnStep) GetName() string {
	return s.Name
}

func (s *ReturnStep) Kind() steptype.Steptype {
	return steptype.Return
}

func (s *ReturnStep) SetOpts(o opts.ExecutionOpts) {
	return
}

func (s *ReturnStep) GetResult() Result {
	return Result{}
}

func (s *ReturnStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) error {
	return nil
}

func (s *ReturnStep) GetCalledTargetName() string {
	return "return"
}

func (s *ReturnStep) GetOpts() opts.ExecutionOpts {
	return opts.ExecutionOpts{}
}

func (s *ReturnStep) SetCalledTarget(t Target) {
	panic("calling target in return")
}
