package step

import (
	"fmt"

	"github.com/kassybas/tame/types/steptype"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/settings"
	"github.com/sirupsen/logrus"
)

type Param struct {
	Name         string
	HasDefault   bool
	DefaultValue interface{}
}
type Target struct {
	GlobalSettings *settings.Settings
	Name           string
	Steps          []Step
	Params         []Param
	Opts           opts.ExecutionOpts
	Variables      []tvar.TVariable
	Summary        string
	Status         int
}

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
	}
}

func (t Target) Make(ctx tcontext.Context, vt vartable.VarTable) ([]interface{}, int, error) {
	vt.AddVariables(ctx.Globals)
	vt, err := resolveParams(vt, t.Params)
	if err != nil {
		return nil, 0, fmt.Errorf("could not resolve parameters in target: %s\n\t%s", t.Name, err)
	}
	var returnValues []interface{}
	for _, s := range t.Steps {
		// Opts
		// TODO: straighten out this mess
		s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, t.Opts, s.GetOpts()))
		newCtx := ctx
		newCtx.Settings.GlobalOpts = s.GetOpts()

		// Run
		results, stdstatus, err := s.RunStep(newCtx, vt)
		if err != nil {
			return nil, stdstatus, fmt.Errorf("%s\n[target: %s, step: %s]", err.Error(), t.Name, s.GetName())
		}
		if s.Kind() == steptype.Return {
			// if return step, break execution
			return results, stdstatus, err
		}
		// Check result status
		if s.GetOpts().CanFail == false && stdstatus != 0 {
			logrus.Errorf("execution failed: status %d\n\ttarget: %s", stdstatus, t.Name)
			return nil, stdstatus, nil
		}
		// Only shell step is allowed to have less results
		allowedLessResults := false
		if s.Kind() == steptype.Shell {
			allowedLessResults = true
		}
		vt, err = updateVarsWithResultVariables(vt, s.ResultNames(), results, allowedLessResults)
		if err != nil {
			return nil, stdstatus, fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)
		}
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%s\n\ttarget: %s", err.Error(), t.Name)
	}
	return returnValues, 0, nil
}

func updateVarsWithResultVariables(vt vartable.VarTable, resultVarNames []string, resultValues []interface{}, allowedLessResults bool) (vartable.VarTable, error) {
	if len(resultVarNames) == 0 {
		return vt, nil
	}
	if len(resultVarNames) > len(resultValues) {
		return vt, fmt.Errorf("too many results expected, too little returned: %d > %d", len(resultVarNames), len(resultValues))
	}
	if len(resultVarNames) != len(resultValues) && !allowedLessResults {
		return vt, fmt.Errorf("return and result variables do not match: %d != %d", len(resultVarNames), len(resultValues))
	}
	// append iterates thorugh names, not values
	err := vt.Append(resultVarNames, resultValues)
	return vt, err
}

// TODO: unify variable resolution
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

func (t Target) isParameter(name string) bool {
	for _, p := range t.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
