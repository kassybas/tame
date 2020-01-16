package waitstep

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type IfStep struct {
	basestep.BaseStep
	condition string
	steps     []step.Step
}

func NewWaitStep(stepDef schema.MergedStepSchema) (*IfStep, error) {
	var newStep IfStep
	var err error
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.If, "if")
	return &newStep, err
}

func (s *IfStep) evalCondition(vt *vartable.VarTable) (bool, error) {
	env := vt.GetAllValues()
	program, err := expr.Compile(s.condition, expr.Env(env))
	if err != nil {
		return false, err
	}
	result, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	resBool, isBool := result.(bool)
	if !isBool {
		return false, fmt.Errorf("if condition expression is not bool: %s -> %s ", s.condition, result)
	}
	return resBool, nil
}

func (s *IfStep) IfStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	if res, err := s.evalCondition(vt); err != nil {
		return err
	}
	if res {

	}
	return step.StepStatus{}
}
