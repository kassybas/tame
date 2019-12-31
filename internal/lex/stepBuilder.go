package lex

import (
	"fmt"

	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/types/steptype"
)

func findCalledTarget(name, caller string, targets map[string]target.Target) (target.Target, error) {
	v, exists := targets[name]
	if exists {
		return v, nil
	}
	return target.Target{}, fmt.Errorf("Target not found: '%s' [called by: '%s']", name, caller)
}

func populateSteps(trg *target.Target, targets map[string]target.Target) error {
	for i := range trg.Steps {
		if trg.Steps[i].Kind() == steptype.Call {
			calledTarget, err := findCalledTarget(trg.Steps[i].GetName(), trg.Name, targets)
			if err != nil {
				return err
			}
			err = populateSteps(&calledTarget, targets)
			if err != nil {
				return err
			}
			trg.Steps[i].(*callstep.CallStep).SetCalledTarget(calledTarget)
		}
	}
	return nil
}
