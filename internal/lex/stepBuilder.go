package lex

import (
	"fmt"

	"github.com/kassybas/mate/internal/step"
	"github.com/kassybas/mate/types/steptype"
)

func findCalledTarget(name, caller string, targets map[string]step.Target) (step.Target, error) {
	v, exists := targets[name]
	if exists {
		return v, nil
	}
	return step.Target{}, fmt.Errorf("Target not found: '%s' [called by: '%s']", name, caller)
}

func populateSteps(trg *step.Target, targets map[string]step.Target) error {
	var err error
	for i := range trg.Steps {
		if trg.Steps[i].Kind == steptype.Call {
			trg.Steps[i].CalledTarget, err = findCalledTarget(trg.Steps[i].CalledTargetName, trg.Name, targets)
			if err != nil {
				return err
			}
			err = populateSteps(&trg.Steps[i].CalledTarget, targets)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
