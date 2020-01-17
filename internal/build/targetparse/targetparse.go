package targetparse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/build/stepbuilder"
	"github.com/kassybas/tame/internal/build/stepparse"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/param"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/stepblock"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/schema"
)

func ParseTeafile(tf schema.Tamefile, ctx *tcontext.Context) (map[string]target.Target, error) {
	targets := make(map[string]target.Target)
	for targetKey, targetValue := range tf.Targets {
		trg, err := buildTarget(targetKey, targetValue, ctx)
		if err != nil {
			return targets, err
		}
		targets[targetKey] = trg
	}
	return targets, nil
}

func buildStep(rawStep map[string]interface{}) (step.Step, error) {
	stepDef, err := stepparse.ParseStepSchema(rawStep)
	if err != nil {
		return nil, err
	}
	return stepbuilder.NewStep(stepDef)

}

func buildSteps(stepDefs []map[string]interface{}) (stepblock.StepBlock, error) {
	steps := make([]step.Step, len(stepDefs))
	for i, stepDef := range stepDefs {
		newStep, err := buildStep(stepDef)
		if err != nil {
			return stepblock.StepBlock{}, err
		}
		steps[i] = newStep
	}
	return stepblock.NewStepBlock(steps), nil
}

func buildTarget(targetKey string, targetDef schema.TargetSchema, ctx *tcontext.Context) (target.Target, error) {
	var err error
	newTarget := target.Target{
		Name: targetKey,
		Ctx:  ctx,
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

func buildParameters(paramDefs map[string]interface{}) ([]param.Param, error) {
	params := []param.Param{}

	for paramKey, paramValue := range paramDefs {
		if !strings.HasPrefix(paramKey, keywords.PrefixReference) {
			return params, fmt.Errorf("arguments must start with '$' symbol: %s (correct: %s%s)", paramKey, keywords.PrefixReference, paramKey)
		}
		newParam := param.Param{
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
