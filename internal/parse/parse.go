package parse

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/schema"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/step"
	"github.com/kassybas/mate/types/target"
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

func buildParameters(paramDefs map[string]interface{}) ([]target.ParamConfig, error) {
	params := []target.ParamConfig{}

	for paramKey, paramValue := range paramDefs {
		newParam := target.ParamConfig{
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
func buildOptsFromString(optsStr string) (opts.ExecutionOpts, error) {
	optsList := strings.Split(optsStr, keywords.OptsSeparator)

	opts := opts.ExecutionOpts{}
	for _, optStr := range optsList {
		opts.Silent = optStr == keywords.OptSilent
		// TODOb: handle all opts
	}
	return opts, nil
}

func buildArguments(interface{}) {
	// TODO: continue here
}

func buildStep(stepDef map[string]interface{}) (step.Step, error) {
	var newStep step.Step
	var err error

	for stepKey, stepValue := range stepDef {
		if !strings.HasPrefix(stepKey, keywords.PrefixTameKeyword) {
			if newStep.Kind == step.Unset {
				newStep.Name = stepKey
				newStep.Kind = step.Call
				newStep.Arguments, err = buildArguments(stepValue)
				if err != nil {
					return newStep, err
				}
				continue
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
				newStep.ResultVars.StdoutVar = stepValue.(string)
			}
		case keywords.ErrVar:
			{
				newStep.ResultVars.StderrVar = stepValue.(string)
			}
		case keywords.RcVar:
			{
				newStep.ResultVars.StdrcVar = stepValue.(string)
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

func buildTarget(targetKey string, targetContainer schema.TargetContainer) (target.Target, error) {
	var err error
	newTarget := target.Target{
		Name: targetKey,
	}

	// Parameters
	newTarget.Params, err = buildParameters(targetContainer.ArgContainer)

	// Steps
	newTarget.Steps, err = buildSteps(targetContainer.BodyContainer)
	if err != nil {
		logrus.Error("Error when parsing targets", err)
	}
	fmt.Printf("%+v", newTarget.Steps)
	panic("OK")

	// for _, depValue := range targetContainer.StepContainer {
	// 	dep, err := buildDependencyDefinition(depValue)
	// 	if err != nil {
	// 		return newTarget, err
	// 	}
	// 	newTarget.Deps = append(newTarget.Deps, dep)
	// }
	return newTarget, err
}
