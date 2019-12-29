package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step/varstep"
)

func buildVarStep(stepDef map[string]interface{}) (varstep.VarStep, error) {
	var newStep varstep.VarStep
	if len(stepDef) > 1 {
		return newStep, fmt.Errorf("multiple variables defined in step, only one allowed: %v", stepDef)
	}
	for k, v := range stepDef {
		if strings.HasPrefix(k, keywords.StepVar) {
			varName, err := parseVariableName(k)
			if err != nil {
				return newStep, err
			}
			err = validateVariableName(varName)
			newStep.Name = varName
			if err != nil {
				return newStep, err
			}
			newStep.Definition = v
			continue
		}
		return varstep.VarStep{}, fmt.Errorf("unknown field in var step: %v", k)
	}
	return newStep, nil
}

func parseVariableName(k string) (string, error) {
	fields := strings.Fields(k)
	if len(fields) > 2 {
		return "", fmt.Errorf("'%s': variable name contains whitespaces", k)
	}
	if len(fields) == 1 {
		return "", fmt.Errorf("'%s': no variable target name found: (correct: var $varname: value)", k)
	}
	return fields[1], nil
}
