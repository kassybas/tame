package step

import (
	"fmt"
	"strconv"

	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/settings"
	"github.com/sirupsen/logrus"
)

type Param struct {
	Name         string
	HasDefault   bool
	DefaultValue string
}
type Target struct {
	GlobalSettings *settings.Settings

	Name      string
	Return    []string
	Steps     []StepI
	Params    []Param
	Opts      opts.ExecutionOpts
	Variables []tvar.Variable
	Summary   string
}

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
	}
}

func (t Target) Run(ctx tcontext.Context, args []tvar.Variable) ([]string, int, error) {
	variables, err := CreateVariables(ctx.Globals, args, t.Params)
	if err != nil {
		return nil, 0, err
	}
	for _, s := range t.Steps {
		// Opts
		// TODO: straighten out this mess
		s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, t.Opts, s.GetOpts()))
		newCtx := ctx
		newCtx.Settings.GlobalOpts = s.GetOpts()

		// Run
		err = s.RunStep(newCtx, variables)
		if err != nil {
			return nil, 0, fmt.Errorf("%s\n\tin target: %s, calling: %s", err.Error(), t.Name, s.GetCalledTargetName())
		}
		// Check result status
		if s.GetOpts().CanFail == false {
			if s.GetResult().StdrcValue != 0 {
				logrus.Errorf("execution failed: status %d\n\tin target:%s", s.GetResult().StdrcValue, t.Name)
				return nil, s.GetResult().StdrcValue, nil
			}
		}
		// Save result variables
		if s.GetResult().ResultVars != nil {
			if len(s.GetResult().ResultValues) != len(s.GetResult().ResultVars) {
				return nil, 0, fmt.Errorf("mismatched number of return and result variables:\n\treturn: %d, result: %d", len(s.GetResult().ResultValues), len(s.GetResult().ResultVars))
			}
		}
		variables = UpdateResultVariables(variables, s.GetResult())
	}

	returnValues, err := createReturnValues(variables, t.Return, t.Name)

	return returnValues, 0, err
}
func UpdateResultVariables(variables map[string]tvar.Variable, r Result) map[string]tvar.Variable {
	if r.StdoutVar != "" {
		variables[r.StdoutVar] = tvar.Variable{Name: r.StdoutVar, Value: r.StdoutValue}
	}
	if r.StderrVar != "" {
		variables[r.StderrVar] = tvar.Variable{Name: r.StderrVar, Value: r.StderrValue}
	}
	if r.StdrcVar != "" {
		variables[r.StdrcVar] = tvar.Variable{Name: r.StdrcVar, Value: strconv.Itoa(r.StdrcValue)}
	}
	for i, v := range r.ResultVars {
		variables[v] = tvar.Variable{Name: v, Value: r.ResultValues[i]}
	}
	return variables
}
