package parse

import (
	"fmt"

	"github.com/kassybas/mate/schema"
	"github.com/kassybas/mate/types/step"
)

// func buildSettings(tfs schema.SetConfig) (settings.Settings, error) {
// 	var settings settings.Settings
// 	settings.UsedShell = tfs.Shell
// 	settings.InitScript = tfs.Init
// 	settings.ShieldEnv = tfs.ShieldEnv
// 	if tfs.DefaultOpts == keywords.OptsNotSet {
// 		settings.DefaultOpts = keywords.OptsDefaultValues
// 	} else {
// 		settings.DefaultOpts = strings.Split(tfs.DefaultOpts, keywords.OptsSeparator)
// 	}
// 	return settings, nil
// }

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
func buildStep(stepDef schema.StepContainer) (step.Step, error) {
	var newStep step.Step
	var err error

	if stepDef.Call == nil && stepDef.Shell == "" {
		return newStep, fmt.Errorf("invalid step configuration: no call or shell defined")
	}
	if stepDef.Call != nil {
		err = populateCallStep(&newStep, stepDef)
		if err != nil {
			return newStep, err
		}
	}
	if stepDef.Shell != "" {
		err = populateShellStep(&newStep, stepDef)
		if err != nil {
			return newStep, err
		}
	}
	newStep.Opts, err = buildOpts(stepDef.Opts)
	return newStep, err
}

func buildSteps(stepDefs []schema.StepContainer) ([]step.Step, error) {
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

func buildTarget(targetKey string, targetContainer schema.TargetContainer) (step.Target, error) {
	var err error
	newTarget := step.Target{
		Name: targetKey,
	}

	// Parameters
	newTarget.Params, err = buildParameters(targetContainer.ArgContainer)
	if err != nil {

		return newTarget, fmt.Errorf("failed to parse steps for '%s'\n\t%s", targetKey, err)
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
