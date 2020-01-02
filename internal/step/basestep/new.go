package basestep

import (
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

func NewBaseStep(stepDef schema.MergedStepSchema, stepType steptype.Steptype, name string) (BaseStep, error) {
	var newStep BaseStep
	var err error
	newStep.kind = stepType
	// Opts
	if stepDef.Opts != nil {
		newStep.opts, err = helpers.BuildOpts(*stepDef.Opts)
		if err != nil {
			return newStep, err
		}
	}
	// Results
	if stepDef.ResultContainers != nil {
		newStep.resultNames = *stepDef.ResultContainers
	}
	// For Loop
	if stepDef.ForLoop.Iterator != "" {
		newStep.iteratorName = stepDef.ForLoop.Iterator
	}
	if stepDef.ForLoop.Iterable != "" {
		newStep.iterableName = stepDef.ForLoop.Iterable
	}
	// If condition
	if stepDef.Condition != nil {
		newStep.ifCondition = *stepDef.Condition
	}
	newStep.name = name
	return newStep, err
}
