package basestep

import (
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type BaseStep struct {
	name        string
	id          *string
	kind        steptype.Steptype
	resultNames []string
	opts        opts.ExecutionOpts
	iteratorVar tvar.TVariable
}

func (s *BaseStep) Kind() steptype.Steptype {
	return s.kind
}
func (s *BaseStep) SetOpts(o opts.ExecutionOpts) {
	s.opts = o
}

func (s *BaseStep) ResultNames() []string {
	return s.resultNames
}

func (s *BaseStep) GetOpts() opts.ExecutionOpts {
	return s.opts
}

func (s *BaseStep) SetIteratorVar(v tvar.TVariable) {
	s.iteratorVar = v
}
func (s *BaseStep) GetIteratorVar() tvar.TVariable {
	return s.iteratorVar
}

func (s *BaseStep) GetName() string {
	return s.name
}
