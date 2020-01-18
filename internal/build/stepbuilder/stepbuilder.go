package stepbuilder

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/step/dumpstep"
	"github.com/kassybas/tame/internal/step/exprstep"
	"github.com/kassybas/tame/internal/step/forstep"
	"github.com/kassybas/tame/internal/step/ifstep"
	"github.com/kassybas/tame/internal/step/returnstep"
	"github.com/kassybas/tame/internal/step/shellstep"
	"github.com/kassybas/tame/internal/step/varstep"
	"github.com/kassybas/tame/internal/step/waitstep"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

func buildSubSteps(subStepDefs []schema.MergedStepSchema) ([]step.Step, error) {
	if subStepDefs == nil {
		return []step.Step{}, nil
	}
	var err error
	steps := make([]step.Step, len(subStepDefs))
	for i := range subStepDefs {
		steps[i], err = NewStep(subStepDefs[i])
		if err != nil {
			return nil, err
		}
	}
	return steps, nil
}

func NewStep(stepDef schema.MergedStepSchema) (step.Step, error) {
	switch stepDef.StepType {
	case steptype.Call:
		{
			return callstep.NewCallStep(stepDef)
		}
	case steptype.Shell:
		{
			return shellstep.NewShellStep(stepDef)
		}
	case steptype.Var:
		{
			return varstep.NewVarStep(stepDef)
		}
	case steptype.Return:
		{
			return returnstep.NewReturnStep(stepDef)
		}
	case steptype.Expr:
		{
			return exprstep.NewExprStep(stepDef)
		}
	case steptype.Wait:
		{
			return waitstep.NewWaitStep(stepDef)
		}
	case steptype.If:
		{
			ifSteps, err := buildSubSteps(stepDef.IfSteps)
			if err != nil {
				return nil, fmt.Errorf("error parsing step in if block:\n\t%s", err.Error())
			}
			elseSteps, err := buildSubSteps(stepDef.ElseSteps)
			if err != nil {
				return nil, fmt.Errorf("error parsing step in else block:\n\t%s", err.Error())
			}
			return ifstep.NewIfStep(stepDef, ifSteps, elseSteps)
		}
	case steptype.For:
		{
			forSteps, err := buildSubSteps(stepDef.ForSteps)
			if err != nil {
				return nil, fmt.Errorf("error parsing step in for-do block:\n\t%s", err.Error())
			}
			return forstep.NewForStep(stepDef, forSteps)
		}
	case steptype.Dump:
		{
			return dumpstep.NewDumpStep(stepDef)
		}
	default:
		{
			return nil, fmt.Errorf("unknown step type: %v", stepDef)
		}
	}
}
