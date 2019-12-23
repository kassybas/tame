package step

import (
	"fmt"
	"strconv"

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
}

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
	}
}

func (t Target) Run(ctx tcontext.Context, vt vartable.VarTable) ([]interface{}, int, error) {
	vt.AddVariables(ctx.Globals)
	vt, err := resolveParams(vt, t.Params)
	if err != nil {
		return nil, 0, fmt.Errorf("could not resolve parameters in target: %s\n\t%s", t.Name, err)
	}
	var returnValues []interface{}
	for _, s := range t.Steps {
		if s.Kind() == steptype.Return {
			rs := s.(*ReturnStep)
			returnValues, err = createReturnValues(vt, rs.Return)
			if err != nil {
				return nil, 0, fmt.Errorf("in target: %s\n\tin step: %s\n\t%s", t.Name, s.GetName(), err)
			}
			break
		}
		// Opts
		// TODO: straighten out this mess
		s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, t.Opts, s.GetOpts()))
		newCtx := ctx
		newCtx.Settings.GlobalOpts = s.GetOpts()

		// Run
		err = s.RunStep(newCtx, vt)
		if err != nil {
			return nil, 0, fmt.Errorf("%s\n[target: %s, step: %s]", err.Error(), t.Name, s.GetName())
		}
		// Check result status
		if s.GetOpts().CanFail == false {
			if s.GetResult().StdStatusValue != 0 {
				logrus.Errorf("execution failed: status %d\n\ttarget: %s", s.GetResult().StdStatusValue, t.Name)
				return nil, s.GetResult().StdStatusValue, nil
			}
		}
		vt, err = updateVarsWithResultVariables(vt, s.GetResult())
		if err != nil {
			return nil, 0, fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)
		}
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%s\n\ttarget: %s", err.Error(), t.Name)
	}
	return returnValues, 0, nil
}

func updateVarsWithResultVariables(vt vartable.VarTable, r Result) (vartable.VarTable, error) {
	if r.StdoutVar != "" {
		vt.Add(r.StdoutVar, r.StdoutValue)
	}
	if r.StderrVar != "" {
		vt.Add(r.StderrVar, r.StderrValue)
	}
	if r.StdStatusVar != "" {
		vt.Add(r.StdStatusVar, strconv.Itoa(r.StdStatusValue))
	}
	if r.ResultValues != nil {
		if len(r.ResultValues) != len(r.ResultNames) {
			return vt, fmt.Errorf("return and result variables do not match: %d != %d", len(r.ResultValues), len(r.ResultNames))
		}
		err := vt.Append(r.ResultNames, r.ResultValues)
		if err != nil {
			return vt, err
		}
	}
	return vt, nil
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

func createReturnValues(vt vartable.VarTable, returnDefinitions []string) ([]interface{}, error) {
	rvs := []interface{}{}
	for _, retDef := range returnDefinitions {
		rv, err := vt.ResolveValue(retDef)
		if err != nil {
			return rvs, err
		}
		rvs = append(rvs, rv)
	}
	return rvs, nil
}

func (t Target) isParameter(name string) bool {
	for _, p := range t.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
