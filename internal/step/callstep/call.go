package callstep

import (
	"fmt"

	"github.com/kassybas/tame/internal/param"
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
	calledTargetSet  bool
}

func NewCallStep(stepDef schema.MergedStepSchema) (*CallStep, error) {
	var err error
	var newStep CallStep
	// Called target
	if stepDef.CalledTargetName == "" {
		return &newStep, fmt.Errorf("missing called target name in call step")
	}
	newStep.calledTargetName = stepDef.CalledTargetName
	// Args
	newStep.arguments = []tvar.TVariable{}
	for k, v := range stepDef.CallArgumentsPassed {
		newStep.arguments = append(newStep.arguments, tvar.NewVariable(k, v))
	}
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Call, newStep.calledTargetName)
	return &newStep, err
}

func (s *CallStep) IsCalledTargetSet() bool {
	return s.calledTargetSet
}

func (s *CallStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	newVt, err := createArgsVartable(s.arguments, s.calledTarget, vt)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("invalid arguments in target call:\n\t%s", err.Error())}
	}

	err = param.ResolveParams(newVt, s.calledTarget.Params)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("could not resolve parameters in target call: %s\n\t%s", s.calledTargetName, err)}
	}
	status := s.calledTarget.Make(newVt, s.GetOpts())
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("step: %s\n\t%s", s.GetName(), status.Err.Error())}
	}
	return status
}

func createArgsVartable(argDefs []tvar.TVariable, calledTarget target.Target, vt *vartable.VarTable) (*vartable.VarTable, error) {
	newVt := vartable.NewVarTable()
	for _, arg := range argDefs {
		if arg.Value() == nil {
			return nil, fmt.Errorf("passing empty(null) argument for target is not allowed %s: '%s: %v'", calledTarget.Name, arg.Name(), arg.Value())
		}
		if !calledTarget.IsParameter(arg.Name()) {
			return nil, fmt.Errorf("unknown parameter for target %s: '%s'", calledTarget.Name, arg.Name())
		}
		resolvedValue, err := vt.ResolveValue(arg.Value())
		if err != nil {
			return nil, fmt.Errorf("passed parameter value cannot be resolved: %s", arg.Value())
		}
		newVt.Add(arg.Name(), resolvedValue)
	}
	return newVt, nil
}

func (s *CallStep) SetCalledTarget(t interface{}) {
	s.calledTarget = t.(target.Target)
	s.calledTargetSet = true
}

func (s *CallStep) GetCalledTarget() target.Target {
	return s.calledTarget
}
