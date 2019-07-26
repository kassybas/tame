package step

import (
	"fmt"
	"strconv"

	"github.com/kassybas/mate/internal/vartable"

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
	Return    string
	Steps     []Step
	Params    []Param
	Opts      opts.ExecutionOpts
	Variables []tvar.VariableI
	Summary   string
}

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
	}
}

func (t Target) Run(ctx tcontext.Context, vt vartable.VarTable) (tvar.VariableI, int, error) {
	vt.AddVariables(ctx.Globals)
	vt, err := resolveParams(vt, t.Params)
	if err != nil {
		return nil, 0, fmt.Errorf("could not resolve parameters in target: %s\n\t%s", t.Name, err)
	}
	for _, s := range t.Steps {
		// Opts
		// TODO: straighten out this mess
		s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, t.Opts, s.GetOpts()))
		newCtx := ctx
		newCtx.Settings.GlobalOpts = s.GetOpts()

		// Run
		err = s.RunStep(newCtx, vt)
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
		// if s.GetResult().ResultVars != nil {
		// 	if len(s.GetResult().ResultValues) != len(s.GetResult().ResultVars) {
		// 		return nil, 0, fmt.Errorf("mismatched number of return and result variables:\n\treturn: %d, result: %d", len(s.GetResult().ResultValues), len(s.GetResult().ResultVars))
		// 	}
		// }
		vt = updateResultVariables(vt, s.GetResult())
	}

	returnValues, err := createReturnValues(vt, t.Return)

	if err != nil {
		return returnValues, 0, fmt.Errorf("%s\n\tin target: %s", err.Error(), t.Name)

	}

	return returnValues, 0, err
}
func updateResultVariables(vt vartable.VarTable, r Result) vartable.VarTable {
	if r.StdoutVar != "" {
		vt.Add(r.StdoutVar, r.StdoutValue)
	}
	if r.StderrVar != "" {
		vt.Add(r.StderrVar, r.StderrValue)
	}
	if r.StdrcVar != "" {
		vt.Add(r.StdrcVar, strconv.Itoa(r.StdrcValue))
	}
	vt.Add(r.ResultVars, r.ResultValue)
	return vt
}

func resolveParams(vt vartable.VarTable, params []Param) (vartable.VarTable, error) {
	for _, p := range params {
		if vt.Exists(p.Name) {
			val, err := vt.GetVar(p.Name)
			if err != nil {
				return vt, err
			}
			vt.Add(p.Name, val.Value())
			continue
		}
		if p.HasDefault {
			vt.Add(p.Name, p.DefaultValue)
			continue
		}
		return vt, fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return vt, nil
}

func createReturnValues(vt vartable.VarTable, returnDefinition string) (tvar.VariableI, error) {
	rv, err := vt.ResolveValue(returnDefinition)
	return rv, err
}
