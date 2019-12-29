package varstep

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type VarStep struct {
	Definition  interface{}
	Opts        opts.ExecutionOpts
	Name        string
	IteratorVar string
	IterableVar string
}

func (s VarStep) GetName() string {
	return s.Name
}

func (s *VarStep) Kind() steptype.Steptype {
	return steptype.Var
}

func (s *VarStep) SetOpts(o opts.ExecutionOpts) {
	s.Opts = o
}

func (s *VarStep) ResultNames() []string {
	// in varstep: the name of the step is equal to the var
	return []string{s.Name}
}

func (s *VarStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	// TODO: eval variables
	value, err := vt.ResolveValue(s.Definition)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.Name, err.Error())}
	}
	return step.StepStatus{Results: []interface{}{value}}
}

func (s *VarStep) GetCalledTargetName() string {
	return s.GetName()
}

func (s *VarStep) GetOpts() opts.ExecutionOpts {
	return s.Opts
}

func (s *VarStep) SetCalledTarget(t interface{}) {
	logrus.Fatal("internal error: calling target in var step")

}

func (s *VarStep) GetIteratorVar() string {
	return s.IteratorVar
}

func (s *VarStep) GetIterableVar() string {
	return s.IterableVar
}
