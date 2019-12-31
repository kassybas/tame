package callstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
)

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
