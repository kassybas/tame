package step

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/tvar"
)

func createReturnValues(variables map[string]tvar.Variable, returnVars []string, targetName string) ([]string, error) {
	returnValues := make([]string, len(returnVars))

	for i, retDef := range returnVars {
		if !strings.HasPrefix(retDef, keywords.PrefixReference) {
			// constant values
			returnValues[i] = retDef
			continue
		}
		_, exists := variables[retDef]
		if !exists {
			return nil, fmt.Errorf("return variable does not exist: '%s'\n\tin target: '%s'", retDef, targetName)
		}
		returnValues[i] = variables[retDef].Value
	}
	return returnValues, nil
}

func CreateVariables(globals []tvar.Variable, args []tvar.Variable, params []Param) (map[string]tvar.Variable, error) {
	variables := make(map[string]tvar.Variable)

	for _, g := range globals {
		variables[g.Name] = g
	}

	for _, p := range params {
		if p.HasDefault {
			variables[p.Name] = tvar.Variable{Name: p.Name, Value: p.DefaultValue}
		}
	}
	for _, a := range args {
		variables[a.Name] = a
	}
	// TODO: check to correct matching of arguments and parameters
	// TODO: check for argument nil values
	return variables, nil
}
