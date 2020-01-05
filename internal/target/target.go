package target

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step"

	"github.com/kassybas/tame/types/vartype"

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

func evalConditionExpression(vt *vartable.VarTable, s step.Step) (bool, error) {
	if s.GetCondition() == "" {
		return true, nil
	}
	env := vt.GetAllValues()
	program, err := expr.Compile(s.GetCondition(), expr.Env(env))
	if err != nil {
		return false, err
	}
	result, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	resBool, isBool := result.(bool)
	if !isBool {
		return false, fmt.Errorf("if condition expression is not bool: %s -> %s ", s.GetCondition(), result)
	}
	return resBool, nil
}

func getIterableValues(iterableIf interface{}, vt *vartable.VarTable) ([]interface{}, error) {

	var iterableVal []interface{}
	switch iterableIf := iterableIf.(type) {
	case string:
		{
			iterable, err := vt.GetVar(iterableIf)
			if err != nil {
				return nil, fmt.Errorf("defined iterable cannot be resolved\n\t%s", err.Error())
			}
			if iterable.Type() != vartype.TListType && iterable.Type() != vartype.TMapType {
				return nil, fmt.Errorf("variable %s is not list or map (type: %T)", iterable.Name(), iterable)
			}
			var isList bool
			iterableVal, isList = iterable.Value().([]interface{})
			if !isList {
				iterableMap := iterable.Value().(map[interface{}]interface{})
				iterableVal = []interface{}{}
				for k := range iterableMap {
					iterableVal = append(iterableVal, k)
				}
			}
		}
	case []interface{}:
		{
			iterableVal = iterableIf
		}
	case map[interface{}]interface{}:
		{
			iterableVal = []interface{}{}
			for k := range iterableIf {
				iterableVal = append(iterableVal, k)
			}
		}
	default:
		{
			return nil, fmt.Errorf("unknown iterable")
		}
	}
	return iterableVal, nil
}

func getIters(vt *vartable.VarTable, s step.Step) (string, []interface{}, error) {
	if s.GetIteratorName() == "" && s.GetIterable() == nil {
		// No iterator and iterable -> no for loop, run once
		return "", []interface{}{""}, nil
	}
	// Iterable
	iterableIf := s.GetIterable()
	if iterableIf == nil {
		// nothing to iterate over -> run zero times
		return "", []interface{}{}, nil
	}
	iterableVal, err := getIterableValues(iterableIf, vt)
	if err != nil {
		return "", nil, err
	}
	// Iterator
	// validate iterator name
	if !strings.HasPrefix(s.GetIteratorName(), keywords.PrefixReference) {
		return "", nil, fmt.Errorf("iterator variable wrong format: %s (should be: %s%s)", s.GetIteratorName(), keywords.PrefixReference, s.GetIteratorName())
	}
	return s.GetIteratorName(), iterableVal, nil
}

func (t Target) Make(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	vt.AddVariables(ctx.Globals)
	err := resolveParams(vt, t.Params)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("could not resolve parameters in target: %s\n\t%s", t.Name, err)}
	}
	status := t.runAllSteps(ctx, vt)
	return status
}

func updateVarsWithResultVariables(vt *vartable.VarTable, resultVarNames []string, resultValues []interface{}, allowedLessResults bool) error {
	if len(resultVarNames) == 0 {
		return nil
	}
	if len(resultVarNames) > len(resultValues) {
		return fmt.Errorf("too many results expected, too little returned: %d > %d", len(resultVarNames), len(resultValues))
	}
	if len(resultVarNames) != len(resultValues) && !allowedLessResults {
		return fmt.Errorf("return and result variables do not match: %d != %d", len(resultVarNames), len(resultValues))
	}

	err := vt.Append(resultVarNames, resultValues)
	return err
}

func resolveParams(vt *vartable.VarTable, params []Param) error {
	for _, p := range params {
		if vt.Exists(p.Name) {
			val, err := vt.GetVar(p.Name)
			if err != nil {
				return err
			}
			vt.Add(p.Name, val.Value())
			continue
		}
		if p.HasDefault {
			vt.Add(p.Name, p.DefaultValue)
			continue
		}
		return fmt.Errorf("parameter without value or default value: %s", p.Name)
	}
	return nil
}

func (t Target) IsParameter(name string) bool {
	for _, p := range t.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
