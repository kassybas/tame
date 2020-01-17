package step

import (
	"reflect"

	"github.com/kassybas/tame/internal/tvar"
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
	SetIteratorVar(tvar.TVariable)
	GetIteratorVar() tvar.TVariable
}

func Clone(s Step) Step {
	indirect := reflect.Indirect(reflect.ValueOf(s))
	newIndirect := reflect.New(indirect.Type())
	newIndirect.Elem().Set(reflect.ValueOf(indirect.Interface()))
	newNamed := newIndirect.Interface()
	casted := newNamed.(Step)
	return casted
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
