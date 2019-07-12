package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/step"
)

func buildParameters(paramDefs map[string]interface{}) ([]step.Param, error) {
	params := []step.Param{}

	for paramKey, paramValue := range paramDefs {
		if !strings.HasPrefix(paramKey, keywords.PrefixReference) {
			return params, fmt.Errorf("arguments must start with '$' symbol: %s (correct: %s%s)", paramKey, keywords.PrefixReference, paramKey)
		}
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

func buildArguments(argDefs map[string]string) ([]step.Variable, error) {
	args := []step.Variable{}
	for argKey, argValue := range argDefs {
		if !strings.HasPrefix(argKey, keywords.PrefixReference) {
			return args, fmt.Errorf("arguments must start with '$' symbol: %s (correct: %s%s)", argKey, keywords.PrefixReference, argKey)
		}
		newArg := step.Variable{
			Name:  argKey,
			Value: argValue,
		}
		args = append(args, newArg)
	}
	return args, nil
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
