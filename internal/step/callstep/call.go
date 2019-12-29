package callstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

// TODO: constructor and make these private
type CallStep struct {
	Name             string
	Arguments        []tvar.TVariable
	Opts             opts.ExecutionOpts
	CalledTargetName string
	CalledTarget     target.Target
	Results          []string
	IteratorVar      string
	IterableVar      string
}

func (s CallStep) GetName() string {
	return s.Name
}

func (s *CallStep) Kind() steptype.Steptype {
	return steptype.Call
}

func (s *CallStep) SetOpts(o opts.ExecutionOpts) {
	s.Opts = o
}

func (s *CallStep) ResultNames() []string {
	return s.Results
}

func (s *CallStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	// TODOb: resolve global variables too
	args, err := createArgsVartable(s.Arguments, s.CalledTarget, vt)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.Name, err.Error())}
	}
	status := s.CalledTarget.Make(ctx, args)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.Name, status.Err.Error())}
	}
	return status
}

func createArgsVartable(argDefs []tvar.TVariable, calledTarget target.Target, vt vartable.VarTable) (vartable.VarTable, error) {
	argsVarTable := vartable.NewVarTable()
	for _, arg := range argDefs {
		if !calledTarget.IsParameter(arg.Name()) {
			return argsVarTable, fmt.Errorf("unknown parameter for target %s: '%s'", calledTarget.Name, arg.Name())
		}
		if arg.Value() == nil {
			return argsVarTable, fmt.Errorf("passing empty(null) argument for target %s: '%s: %v'", calledTarget.Name, arg.Name(), arg.Value())
		}
		val, err := vt.ResolveValue(arg.Value())
		if err != nil {
			return argsVarTable, err
		}
		argsVarTable.Add(arg.Name(), val)
	}

	return argsVarTable, nil
}

func (s *CallStep) GetCalledTargetName() string {
	return s.CalledTargetName
}

func (s *CallStep) GetOpts() opts.ExecutionOpts {
	return s.Opts
}

func (s *CallStep) SetCalledTarget(t interface{}) {
	s.CalledTarget = t.(target.Target)
}

func (s *CallStep) GetIteratorVar() string {
	return s.IteratorVar
}

func (s *CallStep) GetIterableVar() string {
	return s.IterableVar
}
