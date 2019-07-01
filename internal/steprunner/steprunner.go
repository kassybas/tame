package steprunner

import (
	"strconv"

	"github.com/kassybas/mate/internal/helpers"
	"github.com/kassybas/mate/types/step"
	"github.com/kassybas/shell-exec/exec"
)

func CreateVariables(globals []step.Variable, args []step.Variable, params []step.Param) (map[string]step.Variable, error) {
	// Should we resolve here?
	variables := make(map[string]step.Variable)

	for _, g := range globals {
		variables[g.Name] = g
	}

	for _, p := range params {
		if p.HasDefault {
			variables[p.Name] = step.Variable{Name: p.Name, Value: p.DefaultValue}
		}
	}
	for _, a := range args {
		variables[a.Name] = a
	}
	// TODO: check to correct matching of arguments and parameters
	// TODO: check for argument nil values
	return variables, nil
}

func ExecuteScript(s step.Step, vars map[string]step.Variable) (step.Result, error) {
	var err error
	opts := exec.Options{}
	envVars := helpers.FormatEnvVars(vars)
	s.Results.StdoutValue, s.Results.StderrValue, s.Results.StdrcValue, err = exec.ShellExec(s.Script, envVars, opts)
	return s.Results, err
}

func UpdateResultVariables(variables map[string]step.Variable, s step.Step) map[string]step.Variable {
	if s.Results.StdoutVar != "" {
		variables[s.Results.StdoutVar] = step.Variable{Name: s.Results.StdoutVar, Value: s.Results.StdoutValue}
	}
	if s.Results.StderrVar != "" {
		variables[s.Results.StderrVar] = step.Variable{Name: s.Results.StderrVar, Value: s.Results.StderrValue}
	}
	if s.Results.StdrcVar != "" {
		variables[s.Results.StdrcVar] = step.Variable{Name: s.Results.StdrcVar, Value: strconv.Itoa(s.Results.StdrcValue)}
	}
	return variables
}

func (c Context) Exec(target step.Target, args []step.Variable) error {
	variables, err := CreateVariables(c.Globals, args, target.Params)
	if err != nil {
		return err
	}
	for _, s := range target.Steps {
		if s.Kind == step.Exec {
			s.Results, err = ExecuteScript(s, variables)
			if err != nil {
				return err
			}
			variables = UpdateResultVariables(variables, s)
		}
		if s.Kind == step.Call {
			err = c.Exec(s.CalledTarget, s.Arguments)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
