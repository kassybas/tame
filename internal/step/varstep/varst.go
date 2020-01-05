package varstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type VarStep struct {
	basestep.BaseStep
	definition interface{}
}

func NewVarStep(stepDef schema.MergedStepSchema) (*VarStep, error) {
	var err error
	var newStep VarStep
	newStep.definition = stepDef.VarValue
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Var, *stepDef.VarName)
	return &newStep, err
}

func (s *VarStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	// TODO: eval variables
	value, err := vt.ResolveValue(s.definition)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.GetName(), err.Error())}
	}
	return step.StepStatus{Results: []interface{}{value}}
}
