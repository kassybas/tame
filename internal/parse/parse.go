package parse

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/schema"
	"github.com/kassybas/mate/types/opts"
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

func buildParameters(paramDefs map[string]interface{}) ([]step.Param, error) {
	params := []step.Param{}

	for paramKey, paramValue := range paramDefs {
		newParam := step.Param{
			Name: paramKey,
		}
		switch paramValue.(type) {
		case string:
			{
				newParam.HasDefault = true
				newParam.DefaultValue = paramValue.(string)
			}
		default:
			// nil or unknown type
			newParam.HasDefault = false
		}
		params = append(params, newParam)
	}
	return params, nil
}
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
func buildOpts(optsDef []string) (opts.ExecutionOpts, error) {
	opts := opts.ExecutionOpts{}
	for _, opt := range optsDef {
		if opt == keywords.OptSilent {
			opts.Silent = true
		}
		// TODOb: handle all opts
	}
	return opts, nil
}

func buildArguments(argDefs map[string]string) ([]step.Variable, error) {
	args := []step.Variable{}
	for argKey, argValue := range argDefs {
		newArg := step.Variable{
			Name:  argKey,
			Value: argValue,
		}
		args = append(args, newArg)
	}
	return args, nil
}

func populateCallStep(newStep *step.Step, stepDef schema.StepContainer) error {
	var err error
	var keys []string
	for key := range stepDef.Call {
		keys = append(keys, key)
	}
	if len(keys) != 1 {
		return fmt.Errorf("multiple calls defined in single step: %s", keys)
	}
	newStep.Kind = step.Call
	newStep.CalledTargetName = keys[0]
	newStep.Results.ResultVars = stepDef.Result
	newStep.Arguments, err = buildArguments(stepDef.Call[newStep.CalledTargetName])
	return err
}

func populateShellStep(newStep *step.Step, stepDef schema.StepContainer) error {
	var err error
	if newStep.Kind != step.Unset {
		return fmt.Errorf("invalid step configuration: no call or shell defined")
	}
	newStep.Kind = step.Exec
	newStep.Script = stepDef.Shell
	if stepDef.Out != "" {
		newStep.Results.StdoutVar = stepDef.Out
	}
	if stepDef.Out != "" {
		newStep.Results.StderrVar = stepDef.Err
	}
	if stepDef.Rc != "" {
		newStep.Results.StderrVar = stepDef.Rc
	}
	return err
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

	// Steps
	newTarget.Steps, err = buildSteps(targetContainer.BodyContainer)

	if err != nil {
		logrus.Error("failed to parse steps for '%s'\n%s", targetKey, err)
	}
	newTarget.Summary = targetContainer.Summary

	newTarget.Return = targetContainer.ReturnContainer
	return newTarget, err
}
