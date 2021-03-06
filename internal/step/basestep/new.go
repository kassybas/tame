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
	newStep.name = name
	newStep.id = stepDef.ID
	return newStep, err
}
