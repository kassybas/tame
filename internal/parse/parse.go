package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"

	"github.com/kassybas/mate/internal/helpers"
	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/step"
	"github.com/kassybas/mate/schema"
)

func ParseTeafile(tf schema.Tamefile) (map[string]step.Target, error) {

	targets := make(map[string]step.Target)
	for targetKey, targetValue := range tf.Targets {
		trg, err := buildTarget(targetKey, targetValue)
		if err != nil {
			return targets, err
		}
		targets[targetKey] = trg
	}
	return targets, nil
}

func determineStepType(stepDef map[string]interface{}) (steptype.Steptype, error) {
	multiDefError := fmt.Errorf("type of step cannot be determined: more than on of: (var|sh|call|return) only one should be defined in each step")
	sType := steptype.Unset
	for k := range stepDef {
		if k == keywords.StepShell {
			if sType != steptype.Unset {
				return sType, multiDefError
			}
			sType = steptype.Shell
			continue
		}
		if strings.HasPrefix(k, keywords.StepVar) {
			if sType != steptype.Unset {
				return sType, multiDefError
			}
			sType = steptype.Var
			continue
		}
		if strings.HasPrefix(k, keywords.StepCall) {
			if sType != steptype.Unset {
				return sType, multiDefError
			}
			sType = steptype.Call
			continue
		}
		if k == keywords.StepReturn {
			if sType != steptype.Unset {
				return sType, multiDefError
			}
			sType = steptype.Return
			continue
		}
	}
	if sType == steptype.Unset {
		return sType, fmt.Errorf("undeterminable step type: must be (var|sh|call|return)")
	}
	return sType, nil
}

func buildStep(stepDef map[string]interface{}) (step.Step, error) {
	stepType, err := determineStepType(stepDef)
	if err != nil {
		return nil, err
	}
	switch stepType {
	case steptype.Call:
		{
			newStep, err := buildCallStep(stepDef)
			return &newStep, err
		}
	case steptype.Shell:
		{
			newStep, err := buildShellStep(stepDef)
			return &newStep, err
		}
	case steptype.Var:
		{
			newStep, err := buildVarStep(stepDef)
			return &newStep, err
		}
	case steptype.Return:
		{
			newStep, err := buildReturnStep(stepDef)
			return &newStep, err
		}
	}

	return nil, fmt.Errorf("internal parsing error: parsing did not finish succesfully")
}

func buildSteps(stepDefs []map[string]interface{}) ([]step.Step, error) {
	steps := []step.Step{}
	for _, stepDef := range stepDefs {
		newStep, err := buildStep(stepDef)
		if err != nil {
			return steps, err
		}
		steps = append(steps, newStep)
	}
	return steps, nil
}

func buildTarget(targetKey string, targetDef schema.TargetDefinition) (step.Target, error) {
	var err error
	newTarget := step.Target{
		Name: targetKey,
	}

	newTarget.Opts, err = helpers.BuildOpts(targetDef.OptsDefinition)
	if err != nil {
		return newTarget, fmt.Errorf("failed to parse opts for '%s'\n\t%s", targetKey, err)
	}

	// Parameters
	newTarget.Params, err = buildParameters(targetDef.ArgDefinition)
	if err != nil {
		return newTarget, fmt.Errorf("failed to parse parameters for '%s'\n\t%s", targetKey, err)
	}
	// Steps
	newTarget.Steps, err = buildSteps(targetDef.StepDefinition)
	if err != nil {
		return newTarget, fmt.Errorf("failed to parse steps for target '%s'\n\t%s", targetKey, err)
	}

	newTarget.Summary = targetDef.Summary

	return newTarget, err
}

func validateVariableName(name string) error {
	if !strings.HasPrefix(name, keywords.PrefixReference) {
		return fmt.Errorf("variables and arguments must start with '$' symbol: %s (correct: %s%s)", name, keywords.PrefixReference, name)
	}
	return nil
}

func parseOpts(v interface{}) (opts.ExecutionOpts, error) {
	var o opts.ExecutionOpts
	var err error
	switch v.(type) {
	case string:
		{
			o, err = helpers.BuildOpts([]string{v.(string)})
			if err != nil {
				return o, err
			}
		}
	case []interface{}:
		{
			optsSlice, err := ifaceSliceToStringSlice(v.([]interface{}))
			if err != nil {
				return o, err
			}
			o, err = helpers.BuildOpts(optsSlice)
			if err != nil {
				return o, err
			}
		}
	default:
		{
			return o, fmt.Errorf("unknown opts: %s (type %T)", v, v)
		}
	}
	return o, nil
}

func ifaceToString(v interface{}) (string, error) {
	switch v.(type) {
	case string:
		{
			return v.(string), nil
		}
	}
	return "", fmt.Errorf("non-string type: %v (type %T)", v, v)
}

func ifaceSliceToStringSlice(v []interface{}) ([]string, error) {
	res := []string{}
	for i := range v {
		switch v[i].(type) {
		case string:
			{
				res = append(res, v[i].(string))
			}
		default:
			return res, fmt.Errorf("non-string type: %v (type %T)", v[i], v[i])
		}
	}
	return res, nil
}

func buildParameters(paramDefs map[string]interface{}) ([]step.Param, error) {
	params := []step.Param{}

	for paramKey, paramValue := range paramDefs {
		if !strings.HasPrefix(paramKey, keywords.PrefixReference) {
			return params, fmt.Errorf("arguments must start with '$' symbol: %s (correct: %s%s)", paramKey, keywords.PrefixReference, paramKey)
		}
		newParam := step.Param{
			Name: paramKey,
		}
		switch paramValue.(type) {
		case nil:
			newParam.HasDefault = false
		default:
			{
				newParam.HasDefault = true
				newParam.DefaultValue = paramValue
			}
		}
		params = append(params, newParam)
	}
	return params, nil
}
