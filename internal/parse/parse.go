package parse

import (
	"strings"

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
func buildOptsFromString(optsStr string) (opts.ExecutionOpts, error) {
	optsList := strings.Split(optsStr, keywords.OptsSeparator)

	opts := opts.ExecutionOpts{}
	for _, optStr := range optsList {
		opts.Silent = optStr == keywords.OptSilent
		// TODOb: handle all opts
	}
	return opts, nil
}

func buildArguments(argDefs interface{}) ([]step.Variable, error) {
	argMap := argDefs.(map[interface{}]interface{})
	args := []step.Variable{}
	for argKey, argValue := range argMap {
		newArg := step.Variable{
			Name:  argKey.(string),
			Value: argValue.(string),
		}
		args = append(args, newArg)
	}
	return args, nil
}

func buildStep(stepDef map[string]interface{}) (step.Step, error) {
	var newStep step.Step
	var err error

	for stepKey, stepValue := range stepDef {
		if !strings.HasPrefix(stepKey, keywords.PrefixTameKeyword) {
			if newStep.Kind == step.Unset {
				newStep.CalledTargetName = stepKey
				newStep.Kind = step.Call
				newStep.Arguments, err = buildArguments(stepValue)
				if err != nil {
					return newStep, err
				}
			} else {
				logrus.Warn("Ignoring called target because step has different kind set (.exec?)")
			}
		}
		switch stepKey {
		case keywords.Opts:
			{
				newStep.Opts, err = buildOptsFromString(stepValue.(string))
				if err != nil {
					return newStep, err
				}
			}
		case keywords.OutVar:
			{
				newStep.Results.StdoutVar = stepValue.(string)
			}
		case keywords.ErrVar:
			{
				newStep.Results.StderrVar = stepValue.(string)
			}
		case keywords.RcVar:
			{
				newStep.Results.StdrcVar = stepValue.(string)
			}
		case keywords.Exec:
			{
				newStep.Kind = step.Exec
				newStep.Script = stepValue.(string)
			}
		}
	}
	return newStep, err
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
		logrus.Error("Error when parsing targets", err)
	}
	newTarget.Summary = targetContainer.Summary

	newTarget.Return = targetContainer.ReturnContainer
	return newTarget, err
}
