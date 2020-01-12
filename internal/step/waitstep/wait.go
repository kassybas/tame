package waitstep

import (
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type WaitStep struct {
	basestep.BaseStep
}

func NewWaitStep(stepDef schema.MergedStepSchema) (*WaitStep, error) {
	var newStep WaitStep
	var err error
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Wait, "wait")
	return &newStep, err
}

func (s *WaitStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	// waiting is handled in target runner
	return step.StepStatus{}
}
