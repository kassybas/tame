package step

import (
	"github.com/kassybas/tame/internal/vartable"

	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type Step interface {
	GetName() string
	Kind() steptype.Steptype
	ResultNames() []string
	GetOpts() opts.ExecutionOpts
	SetOpts(opts.ExecutionOpts)
	RunStep(tcontext.Context, *vartable.VarTable) StepStatus
	GetIteratorName() string
	GetIterable() interface{}
	GetCondition() string
}

type StepStatus struct {
	Results            []interface{}
	ResultNames        []string
	Stdstatus          int
	IsBreaking         bool
	AllowedLessResults bool
	IsSync             bool
	Err                error
}
