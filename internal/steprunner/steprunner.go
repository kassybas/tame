package steprunner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/step"
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/sirupsen/logrus"

	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
	"github.com/kassybas/shell-exec/exec"
)

func CreateVariables(globals []tvar.Variable, args []tvar.Variable, params []step.Param) (map[string]tvar.Variable, error) {
	variables := make(map[string]tvar.Variable)

	for _, g := range globals {
		variables[g.Name] = g
	}

	for _, p := range params {
		if p.HasDefault {
			variables[p.Name] = tvar.Variable{Name: p.Name, Value: p.DefaultValue}
		}
	}
	for _, a := range args {
		variables[a.Name] = a
	}
	// TODO: check to correct matching of arguments and parameters
	// TODO: check for argument nil values
	return variables, nil
}

func ExecuteScript(ctx tcontext.Context, s step.Step, vars map[string]tvar.Variable) (step.Result, error) {
	var err error
	// ignore result if neither stdout variable and stderr variable is defined
	ignoreResult := s.Results.StderrVar == "" && s.Results.StdoutVar == ""
	opts := exec.Options{
		Silent:       s.Opts.Silent,
		ShellPath:    ctx.Settings.UsedShell,
		IgnoreResult: ignoreResult,
	}
	envVars := step.FormatEnvVars(vars)
	prefixedScript := ctx.Settings.InitScript + "\n" + s.Script
	s.Results.StdoutValue, s.Results.StderrValue, s.Results.StdrcValue, err = exec.ShellExec(prefixedScript, envVars, opts)
	return s.Results, err
}

func UpdateResultVariables(variables map[string]tvar.Variable, s step.Step) map[string]tvar.Variable {
	if s.Results.StdoutVar != "" {
		variables[s.Results.StdoutVar] = tvar.Variable{Name: s.Results.StdoutVar, Value: s.Results.StdoutValue}
	}
	if s.Results.StderrVar != "" {
		variables[s.Results.StderrVar] = tvar.Variable{Name: s.Results.StderrVar, Value: s.Results.StderrValue}
	}
	if s.Results.StdrcVar != "" {
		variables[s.Results.StdrcVar] = tvar.Variable{Name: s.Results.StdrcVar, Value: strconv.Itoa(s.Results.StdrcValue)}
	}
	for i, v := range s.Results.ResultVars {
		variables[v] = tvar.Variable{Name: v, Value: s.Results.ResultValues[i]}
	}
	return variables
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

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
	}
}

func Run(c tcontext.Context, target step.Target, args []tvar.Variable) ([]string, int, error) {
	variables, err := CreateVariables(c.Globals, args, target.Params)
	if err != nil {
		return nil, 0, err
	}
	for _, s := range target.Steps {
		// TODO: fix this
		c.Settings.GlobalOpts = mergeOpts(c.Settings.GlobalOpts, target.Opts, s.Opts)
		s.Opts = mergeOpts(c.Settings.GlobalOpts, target.Opts, s.Opts)
		if s.Kind == steptype.Shell {
			s.Results, err = ExecuteScript(c, s, variables)
			if err != nil {
				return nil, 0, err
			}
			if s.Opts.CanFail == false {
				if s.Results.StdrcValue != 0 {
					logrus.Errorf("execution failed: status %d\n\tin target:%s", s.Results.StdrcValue, target.Name)
					return nil, s.Results.StdrcValue, nil
				}
			}
			variables = UpdateResultVariables(variables, s)
		}
		if s.Kind == steptype.Call {
			stepArgs, err := resolveArgs(s.Arguments, variables)
			if err != nil {
				return nil, 0, err
			}
			returnedValues, rc, err := Run(c, s.CalledTarget, stepArgs)
			if s.Opts.CanFail == false {
				if s.Results.StdrcValue != 0 {
					return nil, rc, nil
				}
			}
			if err != nil {
				return nil, 0, err
			}

			if s.Results.ResultVars != nil {
				if len(returnedValues) != len(s.Results.ResultVars) {
					return nil, 0, fmt.Errorf("mismatched number of return and result variables:\n\treturn: %d, result: %d\n\tin target: %s, calling: %s", len(returnedValues), len(s.Results.ResultVars), target.Name, s.CalledTargetName)
				}

				s.Results.ResultValues = make([]string, len(returnedValues))
				for i := range returnedValues {
					s.Results.ResultValues[i] = returnedValues[i]
				}
				variables = UpdateResultVariables(variables, s)
			}
		}
	}
	returnValues, err := createReturnValues(variables, target.Return, target.Name)

	return returnValues, 0, err
}

func createReturnValues(variables map[string]tvar.Variable, returnVars []string, targetName string) ([]string, error) {
	returnValues := make([]string, len(returnVars))

	for i, retDef := range returnVars {
		if !strings.HasPrefix(retDef, keywords.PrefixReference) {
			// constant values
			returnValues[i] = retDef
			continue
		}
		_, exists := variables[retDef]
		if !exists {
			return nil, fmt.Errorf("return variable does not exist: '%s'\n\tin target: '%s'", retDef, targetName)
		}
		returnValues[i] = variables[retDef].Value
	}
	return returnValues, nil
}
