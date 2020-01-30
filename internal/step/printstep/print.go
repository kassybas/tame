package printstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type PrintStep struct {
	basestep.BaseStep
	value interface{}
}

func NewPrintStep(stepDef schema.MergedStepSchema) (*PrintStep, error) {
	var newStep PrintStep
	var err error
	newStep.value = stepDef.Print
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Wait, fmt.Sprintf("print: %v", newStep.value))
	return &newStep, err
}

func (s *PrintStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	sourceVal, err := vt.ResolveValue(s.value)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("source value cannot be resolved print step: %s\n\t%s", s.GetName(), err.Error())}
	}
	ymlVal, err := helpers.GetFormattedValue(sourceVal, "yaml")
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("source value cannot be converted in print: %s\n\t%s", s.GetName(), err.Error())}
	}
	fmt.Print(ymlVal)
	return step.StepStatus{}
}
