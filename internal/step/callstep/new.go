package callstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type CallStep struct {
	basestep.BaseStep
	calledTargetName string
	arguments        []tvar.TVariable
	calledTarget     target.Target
}

func NewCallStep(stepDef schema.MergedStepSchema) (*CallStep, error) {
	var err error
	var newStep CallStep
	// Called target
	if stepDef.CalledTargetName == nil {
		return &newStep, fmt.Errorf("missing called target name in call step")
	}
	newStep.calledTargetName = *stepDef.CalledTargetName
	// Args
	newStep.arguments = []tvar.TVariable{}
	for k, v := range stepDef.CallArgumentsPassed {
		newStep.arguments = append(newStep.arguments, tvar.NewVariable(k, v))
	}
	name := fmt.Sprintf("%s %s", keywords.StepCall, newStep.calledTargetName)
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Call, name)
	return &newStep, err
}
