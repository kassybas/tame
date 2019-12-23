package step

import (
	"fmt"

	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type CallStep struct {
	Name             string
	Arguments        []tvar.TVariable
	Opts             opts.ExecutionOpts
	Results          Result
	CalledTargetName string
	CalledTarget     Target
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

func (s *CallStep) GetResult() Result {
	return s.Results
}
func (s *CallStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) error {
	// TODOb: resolve global variables too
	args, err := createArgsVartable(s.Arguments, s.CalledTarget, vt)
	if err != nil {
		return fmt.Errorf("in step: %s\n\t%s", s.Name, err.Error())
	}
	s.Results.ResultValues, s.Results.StdStatusValue, err = s.CalledTarget.Run(ctx, args)
	if err != nil {
		return fmt.Errorf("in step: %s\n\t%s", s.Name, err.Error())
	}
	if len(s.Results.ResultNames) != 0 && len(s.Results.ResultValues) != len(s.Results.ResultNames) {
		return fmt.Errorf("mismatch count of return values and result variables: %d != %d", len(s.Results.ResultValues), len(s.Results.ResultNames))
	}
	return nil
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
		argVar, err := vt.ResolveVar(arg)
		if err != nil {
			return argsVarTable, err
		}
		argsVarTable.AddVar(argVar)
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
