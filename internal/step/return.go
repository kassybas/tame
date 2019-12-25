package step

import (
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type ReturnStep struct {
	Name      string
	Arguments []tvar.TVariable
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

func (s *ReturnStep) ResultNames() []string {
	return []string{}
}

func (s *ReturnStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) ([]interface{}, int, error) {
	rvs := []interface{}{}
	for _, retDef := range s.Return {
		rv, err := vt.ResolveValue(retDef)
		if err != nil {
			return rvs, 1, err
		}
		rvs = append(rvs, rv)
	}
	return rvs, 0, nil
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
