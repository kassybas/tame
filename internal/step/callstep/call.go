package callstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
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
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Call, newStep.calledTargetName)
	return &newStep, err
}

func (s *CallStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	// TODOb: resolve global variables too
	args, err := createArgsVartable(s.arguments, s.calledTarget, vt)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.GetName(), err.Error())}
	}
	status := s.calledTarget.Make(ctx, args)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.GetName(), status.Err.Error())}
	}
	return status
}

func createArgsVartable(argDefs []tvar.TVariable, calledTarget target.Target, vt vartable.VarTable) (vartable.VarTable, error) {
	argsVarTable := vartable.NewVarTable()
	for _, arg := range argDefs {
		if !calledTarget.IsParameter(arg.Name()) {
			return argsVarTable, fmt.Errorf("unknown parameter for target %s: '%s'", calledTarget.Name, arg.Name())
		}
		if arg.Value() == nil {
			return argsVarTable, fmt.Errorf("passing empty(null) argument for target %s: '%s: %v'", calledTarget.Name, arg.Name(), arg.Value())
		}
		val, err := vt.ResolveValue(arg.Value())
		if err != nil {
			return argsVarTable, err
		}
		argsVarTable.Add(arg.Name(), val)
	}

	return argsVarTable, nil
}

func (s *CallStep) SetCalledTarget(t interface{}) {
	s.calledTarget = t.(target.Target)
}
