package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/types/steptype"
	"github.com/sirupsen/logrus"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/step/exprstep"
	"github.com/kassybas/tame/internal/step/returnstep"
	"github.com/kassybas/tame/internal/step/shellstep"
	"github.com/kassybas/tame/internal/step/varstep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/schema"
)

func ParseTeafile(tf schema.Tamefile) (map[string]target.Target, error) {
	targets := make(map[string]target.Target)
	for targetKey, targetValue := range tf.Targets {
		trg, err := buildTarget(targetKey, targetValue)
		if err != nil {
			return targets, err
		}
		targets[targetKey] = trg
	}
	return targets, nil
}

func buildStep(rawStep map[string]interface{}) (step.Step, error) {
	stepDef, stepType, err := ParseStepSchema(rawStep)
	if err != nil {
		return nil, err
	}
	switch stepType {
	case steptype.Call:
		{
			return callstep.NewCallStep(stepDef)
		}
	case steptype.Shell:
		{
			return shellstep.NewShellStep(stepDef)
		}
	case steptype.Var:
		{
			return varstep.NewVarStep(stepDef)
		}
	case steptype.Return:
		{
			return returnstep.NewReturnStep(stepDef)
		}
	case steptype.Expr:
		{
			return exprstep.NewExprStep(stepDef)
		}
	default:
		{
			logrus.Fatal("unknown step type: ", rawStep)
		}
	}
	return nil, nil
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

func buildTarget(targetKey string, targetDef schema.TargetSchema) (target.Target, error) {
	var err error
	newTarget := target.Target{
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

func buildParameters(paramDefs map[string]interface{}) ([]target.Param, error) {
	params := []target.Param{}

	for paramKey, paramValue := range paramDefs {
		if !strings.HasPrefix(paramKey, keywords.PrefixReference) {
			return params, fmt.Errorf("arguments must start with '$' symbol: %s (correct: %s%s)", paramKey, keywords.PrefixReference, paramKey)
		}
		newParam := target.Param{
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
