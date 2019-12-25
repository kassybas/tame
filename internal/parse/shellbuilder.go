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
		switch k {
		case keywords.StepShell:
			{
				newStep.Script, err = ifaceToString(v)
				if err != nil {
					return newStep, err
				}
			}
		case keywords.ShellErrResult:
			{
				switch v.(type) {
				case string:
					{
						newStep.Results = []string{v.(string)}
					}
				case []interface{}:
					{
						ifaceSliceToStringSlice(v.([]interface{}))
					}
				default:
					{
						return newStep, fmt.Errorf("unknown type in shell step: %s: %s", k, v)
					}
				}
			}
		case keywords.Opts:
			{
				newStep.Opts, err = parseOpts(v)
				if err != nil {
					return newStep, err
				}
			}
		default:
			{
				return newStep, fmt.Errorf("unknown key in shell step: %s: %s", k, v)
			}
		}
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
