package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/tvar"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step/callstep"
)

func parseCalledTargetName(k string) (string, error) {
	fields := strings.Fields(k)
	if len(fields) > 2 {
		return "", fmt.Errorf("'%s': called target name contains whitespaces", k)
	}
	if len(fields) == 1 {
		return "", fmt.Errorf("'%s': no called target name found", k)
	}
	return fields[1], nil
}

func parseCallStepArgs(argDefs map[interface{}]interface{}) ([]tvar.TVariable, error) {
	args := []tvar.TVariable{}
	for argKey, argValue := range argDefs {

		varName, err := ifaceToString(argKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments %s\n\t%s", argKey, err)
		}
		if err := validateVariableName(varName); err != nil {
			return nil, err
		}
		// TODO: hanlde nonstring argKeys
		newArg := tvar.NewVariable(tvar.ConvertKeyToString(argKey), argValue)
		args = append(args, newArg)
	}
	return args, nil
}

func parseCallStepHeader(newStep *callstep.CallStep, header string, value interface{}) error {
	var err error
	newStep.Name = header
	newStep.CalledTargetName, err = parseCalledTargetName(header)
	if err != nil {
		return err
	}
	if value != nil {
		switch value.(type) {
		case map[interface{}]interface{}:
			{
				newStep.Arguments, err = parseCallStepArgs(value.(map[interface{}]interface{}))
				if err != nil {
					return err
				}
			}
		default:
			{
				return fmt.Errorf("unknown argument type in %s: %v (type %T)", header, value, value)
			}
		}
	}
	return nil
}

func parseCallStepResults(value interface{}) ([]string, error) {
	switch value.(type) {
	case string:
		{
			return []string{value.(string)}, nil
		}
	case []interface{}:
		{
			res, err := ifaceSliceToStringSlice(value.([]interface{}))
			if err != nil {
				return nil, fmt.Errorf("failed to parse result variables: %v\n\t%s", value, err)

			}
			return res, nil
		}
	default:
		{
			return nil, fmt.Errorf("unknown result type: %s (type %T)", value, value)
		}
	}
}

func buildCallStep(stepDef map[string]interface{}) (callstep.CallStep, error) {
	var newStep callstep.CallStep
	var err error
	for k, v := range stepDef {
		if strings.HasPrefix(k, keywords.StepCall) {
			if newStep.CalledTargetName != "" {
				return newStep, fmt.Errorf("multiple call defined in a single step: 'call %s' and '%s'", newStep.CalledTargetName, k)
			}
			if err = parseCallStepHeader(&newStep, k, v); err != nil {
				return newStep, err
			}
			continue
		}
		if k == keywords.StepCallResult {
			newStep.Results, err = parseCallStepResults(v)
			if err != nil {
				return newStep, err
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
		if k == keywords.StepFor {
			newStep.IteratorVar, newStep.IterableVar, err = parseForLoop(v)
			continue
		}
		return newStep, fmt.Errorf("unknown field in call step: %s", k)
	}
	return newStep, nil
}
