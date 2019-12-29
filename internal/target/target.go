package target

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step"

	"github.com/kassybas/tame/types/steptype"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/settings"
)

type Param struct {
	Name         string
	HasDefault   bool
	DefaultValue interface{}
}

type Target struct {
	GlobalSettings *settings.Settings
	Name           string
	Steps          []step.Step
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

func (t Target) runStep(s step.Step, ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	// Opts
	// TODO: straighten out this mess
	s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, t.Opts, s.GetOpts()))
	newCtx := ctx
	newCtx.Settings.GlobalOpts = s.GetOpts()
	// Run
	status := s.RunStep(ctx, vt)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("[target: %s]:: %s", t.Name, status.Err.Error())}
	}
	// Breaking if it was breaking (return step) or the called step exec failed with non-zero exit
	status.IsBreaking = status.IsBreaking || (s.GetOpts().CanFail == false && status.Stdstatus != 0)
	return status
}

func getIters(vt vartable.VarTable, s step.Step) (string, []tvar.TVariable, error) {
	if s.GetIterableVar() == "" {
		return "", nil, nil
	}
	if !strings.HasPrefix(s.GetIteratorVar(), keywords.PrefixReference) {
		return "", nil, fmt.Errorf("iterator variable wrong format: %s (should be: %s%s)", s.GetIteratorVar(), keywords.PrefixReference, s.GetIteratorVar())
	}
	iterable, err := vt.GetVar(s.GetIterableVar())
	v, isList := iterable.(tvar.TList)
	if !isList {
		return "", nil, fmt.Errorf("iterable variable %s is not list (type: %T)", iterable.Name(), iterable)
	}
	return s.GetIteratorVar(), v.Value().([]tvar.TVariable), err
}

func (t Target) Make(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	vt.AddVariables(ctx.Globals)
	vt, err := resolveParams(vt, t.Params)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("could not resolve parameters in target: %s\n\t%s", t.Name, err)}
	}
	for _, s := range t.Steps {
		// TODO: refactor to more dry
		if s.GetIterableVar() == "" && s.GetIteratorVar() == "" {
			status := t.runStep(s, ctx, vt)
			if status.IsBreaking {
				// setting the false so caller does not break
				status.IsBreaking = false
				return status
			}
			vt, err = updateVarsWithResultVariables(vt, s.ResultNames(), status.Results, s.Kind() == steptype.Shell)
			if err != nil {
				return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)}
			}
		} else {
			iterator, iterable, err := getIters(vt, s)
			for _, itVar := range iterable {
				vt.Add(iterator, itVar.Value())
				status := t.runStep(s, ctx, vt)
				if status.IsBreaking {
					// setting the false so does not break
					status.IsBreaking = false
					return status
				}
				vt, err = updateVarsWithResultVariables(vt, s.ResultNames(), status.Results, s.Kind() == steptype.Shell)
				if err != nil {
					return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)}
				}
			}
		}
	}
	return step.StepStatus{}
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

func (t Target) IsParameter(name string) bool {
	for _, p := range t.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
