package parse

import (
	"fmt"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step"
)

func buildShellStep(stepDef map[string]interface{}) (step.ShellStep, error) {
	var err error
	var newStep step.ShellStep
	newStep.Name = keywords.StepShell
	for k, v := range stepDef {
		if k == keywords.StepShell {
			newStep.Script, err = ifaceToString(v)
			if err != nil {
				return newStep, err
			}
			continue
		}
		if k == keywords.ShellOutResult {
			newStep.Results.StdoutVar, err = getVarNameFromIface(v)
			if err != nil {
				return newStep, fmt.Errorf("failed to parse step variable: '%s: %v'\n\t%s", k, v, err)
			}
			continue
		}
		if k == keywords.ShellErrResult {
			newStep.Results.StderrVar, err = getVarNameFromIface(v)
			if err != nil {
				return newStep, fmt.Errorf("failed to parse step variable: '%s: %v'\n\t%s", k, v, err)
			}
			continue
		}
		if k == keywords.ShellStatusResult {
			newStep.Results.StdStatusVar, err = getVarNameFromIface(v)
			if err != nil {
				return newStep, fmt.Errorf("failed to parse step variable: '%s: %v'\n\t%s", k, v, err)
			}
			continue
		}
		if k == keywords.Opts {
			newStep.Opts, err = parseOpts(v)
			if err != nil {
				return newStep, err
			}
			continue
		}
		return newStep, fmt.Errorf("unknown key in shell step: %s: %s", k, v)
	}
	return newStep, nil
}

func getVarNameFromIface(v interface{}) (string, error) {
	varName, err := ifaceToString(v)
	if err != nil {
		return "", err
	}
	err = validateVariableName(varName)
	if err != nil {
		return "", err
	}
	return varName, nil
}
