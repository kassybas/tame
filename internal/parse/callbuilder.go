package parse

import (
	"fmt"
	"strings"
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

func parseCallStepArgs(argDefs interface{}) (map[string]interface{}, error) {
	argMap, ok := argDefs.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("call step must have map as value, got: %T", argDefs)
	}
	args := make(map[string]interface{}, len(argMap))
	for argKey, argValue := range argMap {
		argName, ok := argKey.(string)
		if !ok {
			fmt.Errorf("non-string argument variable name: %v (type %T)", argKey, argKey)
		}
		if err := validateVariableName(argName); err != nil {
			return nil, err
		}
		args[argName] = argValue
	}
	return args, nil
}

// func parseCallStepHeader(newStep *callstep.CallStep, header string, value interface{}) error {
// 	var err error
// 	newStep.Name = header
// 	newStep.CalledTargetName, err = parseCalledTargetName(header)
// 	if err != nil {
// 		return err
// 	}
// 	if value != nil {
// 		switch value.(type) {
// 		case map[interface{}]interface{}:
// 			{
// 				// newStep.Arguments, err = parseCallStepArgs(value.(map[interface{}]interface{}))
// 				// if err != nil {
// 				// 	return err
// 				// }
// 			}
// 		default:
// 			{
// 				return fmt.Errorf("unknown argument type in %s: %v (type %T)", header, value, value)
// 			}
// 		}
// 	}
// 	return nil
// }

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

// func buildCallStep(stepDef map[string]interface{}) (callstep.CallStep, error) {
// 	var newStep callstep.CallStep
// 	var err error
// 	for k, v := range stepDef {
// 		if strings.HasPrefix(k, keywords.StepCall) {
// 			if newStep.CalledTargetName != "" {
// 				return newStep, fmt.Errorf("multiple call defined in a single step: 'call %s' and '%s'", newStep.CalledTargetName, k)
// 			}
// 			if err = parseCallStepHeader(&newStep, k, v); err != nil {
// 				return newStep, err
// 			}
// 			continue
// 		}
// 		if k == keywords.StepCallResult {
// 			newStep.Results, err = parseCallStepResults(v)
// 			if err != nil {
// 				return newStep, err
// 			}
// 			continue
// 		}
// 		if k == keywords.Opts {
// 			newStep.Opts, err = parseOpts(v)
// 			if err != nil {
// 				return newStep, err
// 			}
// 			continue
// 		}
// 		if k == keywords.StepFor {
// 			newStep.IteratorVar, newStep.IterableVar, err = parseForLoop(v)
// 			continue
// 		}
// 		return newStep, fmt.Errorf("unknown field in call step: %s", k)
// 	}
// 	return newStep, nil
// }
