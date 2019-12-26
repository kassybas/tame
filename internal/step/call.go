package step

import (
	"fmt"

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
	CalledTarget     Target
	Results          []string
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

func (s *CallStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) ([]interface{}, int, error) {
	// TODOb: resolve global variables too
	args, err := createArgsVartable(s.Arguments, s.CalledTarget, vt)
	if err != nil {
		return nil, 0, fmt.Errorf("step: %s\n\t%s", s.Name, err.Error())
	}
	resultValues, stdstatus, err := s.CalledTarget.Make(ctx, args)
	if err != nil {
		return resultValues, stdstatus, fmt.Errorf("step: %s\n\t%s", s.Name, err.Error())
	}
	return resultValues, stdstatus, nil
}

func createArgsVartable(argDefs []tvar.TVariable, calledTarget Target, vt vartable.VarTable) (vartable.VarTable, error) {
	argsVarTable := vartable.NewVarTable()
	for _, arg := range argDefs {
		if !calledTarget.isParameter(arg.Name()) {
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

func (s *CallStep) SetCalledTarget(t Target) {
	s.CalledTarget = t
}
