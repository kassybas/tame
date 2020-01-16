package ifstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/steprunner"

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
	ifSteps   []step.Step
	elseSteps []step.Step
}

func NewIfStep(stepDef schema.MergedStepSchema, ifSteps, elseSteps []step.Step) (*IfStep, error) {
	var newStep IfStep
	var err error
	if stepDef.IfCondition == nil {
		return nil, fmt.Errorf("no condition in if step")
	}
	newStep.condition = *stepDef.IfCondition
	newStep.ifSteps = ifSteps
	newStep.elseSteps = elseSteps
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.If, *stepDef.IfCondition)

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

func (s *IfStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	if s.condition == "" {
		return step.StepStatus{Err: fmt.Errorf("empty if condition")}
	}
	res, err := s.evalCondition(vt)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("\n\t while evaluating condition %s:\n\t%s", s.GetName(), err.Error())}
	}
	var status step.StepStatus
	if res {
		status = steprunner.RunAllSteps(s.ifSteps, ctx, vt, s.GetOpts())
	} else {
		status = steprunner.RunAllSteps(s.elseSteps, ctx, vt, s.GetOpts())
	}
	if status.Err != nil {
		status.Err = fmt.Errorf("\n\t in condition %s:\n\t%s", s.GetName(), status.Err.Error())
	}
	return status
}
