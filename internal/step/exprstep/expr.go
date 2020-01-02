package exprstep

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

type ExprStep struct {
	basestep.BaseStep
	expr string
}

func NewExprStep(stepDef schema.MergedStepSchema) (*ExprStep, error) {
	var err error
	var newStep ExprStep
	if stepDef.Expr == nil {
		return &newStep, fmt.Errorf("missing called script in shell step")
	}
	newStep.expr = *stepDef.Expr
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Expr, "expr")
	return &newStep, err
}

func (s *ExprStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	var err error
	env := vt.GetAllValues()
	program, err := expr.Compile(s.expr, expr.Env(env))
	if err != nil {
		return step.StepStatus{Err: err}
	}
	result, err := expr.Run(program, env)
	return step.StepStatus{
		Results:   []interface{}{result},
		Stdstatus: 0,
		Err:       err,
	}
}
