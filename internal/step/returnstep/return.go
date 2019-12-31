package returnstep

import (
	"fmt"
	"log"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type ReturnStep struct {
	Arguments []tvar.TVariable
	Return    []string
}

func (s ReturnStep) GetName() string {
	return "return"
}

func (s *ReturnStep) Kind() steptype.Steptype {
	return steptype.Return
}

func (s *ReturnStep) SetOpts(o opts.ExecutionOpts) {
	return
}

func (s *ReturnStep) ResultNames() []string {
	return []string{}
}

func (s *ReturnStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	rvs := []interface{}{}
	for _, retDef := range s.Return {
		rv, err := vt.ResolveValue(retDef)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("step: %s %v\n\t%s", s.GetName(), s.Return, err.Error())}
		}
		rvs = append(rvs, rv)
	}
	return step.StepStatus{Results: rvs, Stdstatus: 0, Err: nil, IsBreaking: true}
}

func (s *ReturnStep) GetCalledTargetName() string {
	return "return"
}

func (s *ReturnStep) GetOpts() opts.ExecutionOpts {
	return opts.ExecutionOpts{}
}

func (s *ReturnStep) SetCalledTarget(t interface{}) {
	log.Fatal("calling target in return")
}

func (s *ReturnStep) GetIteratorName() string {
	return ""
}

func (s *ReturnStep) GetIterableName() string {
	return ""
}
