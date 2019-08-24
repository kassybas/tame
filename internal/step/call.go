package step

import (
	"fmt"

	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/internal/vartable"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
)

type CallStep struct {
	Name             string
	Arguments        []tvar.VariableI
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
	args, err := createArgsVartable(s.Arguments, vt)
	if err != nil {
		return fmt.Errorf("in step: %s\n\t%s", s.Name, err.Error())
	}
	s.Results.ResultValues, s.Results.StdStatusValue, err = s.CalledTarget.Run(ctx, args)
	if err != nil {
		err = fmt.Errorf("in step: %s\n\t%s", s.Name, err.Error())
	}
	return err
}

func createArgsVartable(argDefs []tvar.VariableI, vt vartable.VarTable) (vartable.VarTable, error) {
	argsVarTable := vartable.NewVarTable()
	for _, arg := range argDefs {
		argVar, err := vt.ResolveVar(arg)
		if err != nil {
			return argsVarTable, err
		}
		argsVarTable.Add(argVar)
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
