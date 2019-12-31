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
