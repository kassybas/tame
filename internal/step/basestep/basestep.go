package basestep

import (
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type BaseStep struct {
	kind         steptype.Steptype
	resultNames  []string
	opts         opts.ExecutionOpts
	iteratorName string
	iterableName string
}

func (s *BaseStep) Kind() steptype.Steptype {
	return steptype.Call
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

func (s *BaseStep) GetIterableName() string {
	return s.iterableName
}
