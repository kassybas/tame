package lex

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/types/steptype"
)

func findCalledTarget(name, caller string, targets map[string]step.Target) (step.Target, error) {
	v, exists := targets[name]
	if exists {
		return v, nil
	}
	return step.Target{}, fmt.Errorf("Target not found: '%s' [called by: '%s']", name, caller)
}

func populateSteps(trg *step.Target, targets map[string]step.Target) error {
	for i := range trg.Steps {
		if trg.Steps[i].Kind() == steptype.Call {
			calledTarget, err := findCalledTarget(trg.Steps[i].GetCalledTargetName(), trg.Name, targets)
			if err != nil {
				return err
			}
			err = populateSteps(&calledTarget, targets)
			if err != nil {
				return err
			}
			trg.Steps[i].SetCalledTarget(calledTarget)
		}
	}
	return nil
}
