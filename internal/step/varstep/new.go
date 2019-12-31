package varstep

import (
	"github.com/kassybas/tame/internal/step/basestep"
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
