package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/step"
	"github.com/kassybas/mate/schema"
)

func populateCallStep(newStep *step.CallStep, stepDef schema.StepDefinition) error {
	var err error
	var keys []string
	for key := range stepDef.Call {
		keys = append(keys, key)
	}
	if len(keys) != 1 {
		return fmt.Errorf("multiple calls defined in single step: %s", keys)
	}
	newStep.CalledTargetName = keys[0]
	newStep.Results.ResultNames = []string{stepDef.Out}
	newStep.Arguments, err = buildArguments(stepDef.Call[newStep.CalledTargetName])
	return err
}

func populateShellStep(newStep *step.ShellStep, stepDef schema.StepDefinition) error {
	var err error
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

func populateVarStep(newStep *step.VarStep, stepDef schema.StepDefinition) error {
	newStep.Definitions = make(map[string]interface{})
	for k, v := range stepDef.Var {
		if !strings.HasPrefix(k, keywords.PrefixReference) {
			return fmt.Errorf("variables must start with '$' symbol: %s (correct: %s%s)", k, keywords.PrefixReference, k)
		}
		newStep.Definitions[k] = v
	}
	return nil
}
