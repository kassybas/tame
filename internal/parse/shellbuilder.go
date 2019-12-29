package parse

import (
	"fmt"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step/shellstep"
)

func buildShellStep(stepDef map[string]interface{}) (shellstep.ShellStep, error) {
	var err error
	var newStep shellstep.ShellStep
	newStep.Name = keywords.StepShell
	for k, v := range stepDef {
		switch k {
		case keywords.StepShell:
			{
				newStep.Script, err = ifaceToString(v)
				if err != nil {
					return newStep, err
				}
				continue
			}
		case keywords.StepFor:
			{
				newStep.IteratorVar, newStep.IterableVar, err = parseForLoop(v)
				continue
			}
		case keywords.StepCallResult:
			{
				switch v.(type) {
				case string:
					{
						newStep.Results = []string{v.(string)}
						continue
					}
				case nil:
					{
						newStep.Results = []string{""}
						continue
					}
				case []interface{}:
					{
						var err error
						newStep.Results, err = ifaceSliceToStringSlice(v.([]interface{}))
						if err != nil {
							return newStep, err
						}
						continue
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
				continue
			}
		default:
			{
				return newStep, fmt.Errorf("unknown key in shell step: %s: %s", k, v)
			}
		}
	}
	return newStep, nil
}

func parseForLoop(forDef interface{}) (string, string, error) {
	f, ok := forDef.(map[interface{}]interface{})
	iterator, ok := f[keywords.StepForIterator]
	if !ok {
		return "", "", fmt.Errorf("missing interator variable in for loop")
	}
	iteratorStr, ok := iterator.(string)
	if !ok {
		return "", "", fmt.Errorf("non-string iterable variable: %v", iterator)
	}
	iterable, ok := f[keywords.StepForIterable]
	if !ok {
		return "", "", fmt.Errorf("missing iterable variable in for loop")
	}
	iterableStr, ok := iterable.(string)
	if !ok {
		return "", "", fmt.Errorf("non-string iterable variable: %v", iterable)
	}
	return iteratorStr, iterableStr, nil
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
