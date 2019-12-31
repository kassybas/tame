package returnstep

import (
	"fmt"
	"log"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
)

func (s *ReturnStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	rvs := []interface{}{}
	for _, retDef := range s.returnNames {
		rv, err := vt.ResolveValue(retDef)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("step: %s %v\n\t%s", s.GetName(), s.returnNames, err.Error())}
		}
		rvs = append(rvs, rv)
	}
	return step.StepStatus{Results: rvs, Stdstatus: 0, Err: nil, IsBreaking: true}
}

func (s *ReturnStep) GetCalledTargetName() string {
	return "return"
}

func (s *ReturnStep) SetCalledTarget(t interface{}) {
	log.Fatal("internal error: calling target in return")
}
