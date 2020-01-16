package basestep

import (
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type BaseStep struct {
	name         string
	kind         steptype.Steptype
	resultNames  []string
	opts         opts.ExecutionOpts
	iteratorName string
	iterable     interface{}
	ifCondition  string
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

func (s *BaseStep) GetIteratorName() string {
	return s.iteratorName
}

func (s *BaseStep) GetIterable() interface{} {
	return s.iterable
}

func (s *BaseStep) GetName() string {
	return s.name
}
