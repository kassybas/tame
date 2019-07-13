package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/step"
	"github.com/kassybas/mate/schema"
	"github.com/kassybas/mate/types/steptype"
)

func populateCallStep(newStep *step.Step, stepDef schema.StepContainer) error {
	var err error
	var keys []string
	for key := range stepDef.Call {
		keys = append(keys, key)
	}
	if len(keys) != 1 {
		return fmt.Errorf("multiple calls defined in single step: %s", keys)
	}
	newStep.Kind = steptype.Call
	newStep.CalledTargetName = keys[0]
	newStep.Results.ResultVars = stepDef.Result
	newStep.Arguments, err = buildArguments(stepDef.Call[newStep.CalledTargetName])
	return err
}

func populateShellStep(newStep *step.Step, stepDef schema.StepContainer) error {
	var err error
	if newStep.Kind != steptype.Unset {
		return fmt.Errorf("invalid step configuration: no call or shell defined")
	}
	newStep.Kind = steptype.Shell
	newStep.Script = stepDef.Shell
	if stepDef.Out != "" {
		if !strings.HasPrefix(stepDef.Out, keywords.PrefixReference) {
			return fmt.Errorf("out variables must start with '$' symbol: %s (correct: %s%s)", stepDef.Out, keywords.PrefixReference, stepDef.Out)
		}
		newStep.Results.StdoutVar = stepDef.Out
	}
	if stepDef.Out != "" {
		newStep.Results.StderrVar = stepDef.Err
		if !strings.HasPrefix(stepDef.Out, keywords.PrefixReference) {
			return fmt.Errorf("err variables must start with '$' symbol: %s (correct: %s%s)", stepDef.Err, keywords.PrefixReference, stepDef.Err)
		}
	}
	if stepDef.Rc != "" {
		newStep.Results.StderrVar = stepDef.Rc
		if !strings.HasPrefix(stepDef.Out, keywords.PrefixReference) {
			return fmt.Errorf("rc variables must start with '$' symbol: %s (correct: %s%s)", stepDef.Rc, keywords.PrefixReference, stepDef.Rc)
		}
	}
	return err
}
