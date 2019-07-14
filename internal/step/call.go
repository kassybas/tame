package step

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
)

type CallStep struct {
	Name             string
	Arguments        []tvar.Variable
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
func (s *CallStep) RunStep(ctx tcontext.Context, vars map[string]tvar.Variable) error {
	// TODOb: resolve global variables too
	args, err := resolveArgs(s.Arguments, vars)
	if err != nil {
		return err
	}
	s.Results.ResultValues, s.Results.StdrcValue, err = s.CalledTarget.Run(ctx, args)
	return nil
}

func resolveArgs(argDefs []tvar.Variable, variables map[string]tvar.Variable) ([]tvar.Variable, error) {
	for i, arg := range argDefs {
		if strings.HasPrefix(arg.Value, keywords.PrefixReference) {
			_, exists := variables[arg.Value]
			if !exists {
				return nil, fmt.Errorf("variable does not exist in context: '%s:%s'", arg.Name, arg.Value)
			}
			argDefs[i].Value = variables[arg.Name].Value
		}
	}
	return argDefs, nil
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
