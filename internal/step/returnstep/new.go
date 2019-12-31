package returnstep

import (
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type ReturnStep struct {
	basestep.BaseStep
	returnNames []string
}

func NewReturnStep(stepDef schema.MergedStepSchema) (*ReturnStep, error) {
	var newStep ReturnStep
	var err error
	if stepDef.Return != nil {
		newStep.returnNames = *stepDef.Return
	} else {
		newStep.returnNames = []string{}
	}
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Return, "return")
	return &newStep, err
}
