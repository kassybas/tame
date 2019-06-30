package parser

import (
	"fmt"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/target"
	"strings"
)

func buildOpts(optsValue string) ([]string, error) {
	return strings.Split(optsValue, keywords.OptsSeparator), nil
}

func buildDepCalledArgs(argValues map[interface{}]interface{}) (map[string]string, error) {
	args := make(map[string]string)
	for argKey, argValue := range argValues {
		args[argKey.(string)] = argValue.(string)
	}
	return args, nil
}

func buildComplexDependencyDefinition(depValue map[interface{}]interface{}) (target.DepConfig, error) {
	var newDep target.DepConfig

	for depKey, depValue := range depValue {
		switch depValue.(type) {
		// Handle opts
		// opts: silent can-fail
		case string:
			{
				if depKey.(string) != keywords.Opts {
					return newDep, fmt.Errorf("unknown keyword: '%s'", depKey)
				}

				var err error
				newDep.Opts, err = buildOpts(depValue.(string))
				if err != nil {
					return newDep, err
				}
			}
		// Handle called arguments
		// example-dep1: { arg1: value1, arg2: value2  }
		case map[interface{}]interface{}:
			{
				newDep.Name = depKey.(string)
				var err error
				newDep.ArgValues, err = buildDepCalledArgs(depValue.(map[interface{}]interface{}))
				if err != nil {
					return newDep, err
				}
			}
		default:
			{
				return newDep, fmt.Errorf("unknown yaml sructure: could not determine dependency type. it should be string 'opts' or map 'target:{arg:value}': %s", depValue)
			}
		}
	}
	return newDep, nil
}

func buildDependencyDefinition(depValue interface{}) (target.DepConfig, error) {

	// Multiple types possible in schema. Eg:
	// target-name:
	//   deps:
	//   - dep1                     // <- no args defined value: string
	//   - dep2: {arg1: "value1"}   // <- arg given: map[string]map[string]string
	switch depValue.(type) {
	case string:
		{
			// Simple dependency
			var newDep target.DepConfig
			newDep.Name = depValue.(string)
			return newDep, nil
		}
	case map[interface{}]interface{}:
		{
			// Complex dependency with arguments and opts
			newDep, err := buildComplexDependencyDefinition(depValue.(map[interface{}]interface{}))
			if err != nil {
				return target.DepConfig{}, err
			}
			return newDep, nil
		}
	}
	return target.DepConfig{}, fmt.Errorf("incorrect yaml sructure: could not determine dependency type: it should be 'list of string' or 'map of string to string': %s", depValue)
}
