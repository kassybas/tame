package parse

import (
	"fmt"

	"github.com/kassybas/mate/internal/helpers"
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
func buildStep(stepDef schema.StepContainer) (step.StepI, error) {
	var err error
	var newStep step.StepI

	if stepDef.Call == nil && stepDef.Shell == "" {
		return nil, fmt.Errorf("invalid step configuration: no call or shell defined")
	}
	// Call
	if stepDef.Call != nil {
		var newCallStep step.CallStep
		err = populateCallStep(&newCallStep, stepDef)
		if err != nil {
			return &newCallStep, err
		}
		newCallStep.Opts, err = helpers.BuildOpts(stepDef.Opts)
		newStep = &newCallStep
	}
	// Shell
	if stepDef.Shell != "" {
		var newShellStep step.ShellStep
		err = populateShellStep(&newShellStep, stepDef)
		if err != nil {
			return &newShellStep, err
		}
		newShellStep.Opts, err = helpers.BuildOpts(stepDef.Opts)
		newStep = &newShellStep
	}
	return newStep, err
}

func buildSteps(stepDefs []schema.StepContainer) ([]step.StepI, error) {
	steps := []step.StepI{}
	for _, stepDef := range stepDefs {
		newStep, err := buildStep(stepDef)
		if err != nil {
			return steps, err
		}
		steps = append(steps, newStep)
	}
	return steps, nil
}

func buildTarget(targetKey string, targetContainer schema.TargetContainer) (step.Target, error) {
	var err error
	newTarget := step.Target{
		Name: targetKey,
	}

	newTarget.Opts, err = helpers.BuildOpts(targetContainer.OptsContainer)
	if err != nil {
		return newTarget, fmt.Errorf("failed to parse opts for '%s'\n\t%s", targetKey, err)
	}

	// Parameters
	newTarget.Params, err = buildParameters(targetContainer.ArgContainer)
	if err != nil {
		return newTarget, fmt.Errorf("failed to parse parameters for '%s'\n\t%s", targetKey, err)
	}

	// Steps
	newTarget.Steps, err = buildSteps(targetContainer.BodyContainer)
	if err != nil {
		return newTarget, fmt.Errorf("failed to parse steps for '%s'\n\t%s", targetKey, err)
	}

	newTarget.Summary = targetContainer.Summary

	newTarget.Return = targetContainer.ReturnContainer
	return newTarget, err
}
