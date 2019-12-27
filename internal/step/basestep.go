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
	RunStep(tcontext.Context, vartable.VarTable) StepStatus
	GetCalledTargetName() string
	SetCalledTarget(Target)
	GetIteratorVar() string
	GetIterableVar() string
}

type StepStatus struct {
	Results    []interface{}
	IsBreaking bool
	Stdstatus  int
	Err        error
}

// type BaseStep struct {
// 	name    string
// 	kind    steptype.Steptype
// 	results []string
// 	opts    opts.ExecutionOpts
// }

// TODO: cleanup
// type Result struct {
// 	StdoutVar      string
// 	StdoutValue    string
// 	StderrVar      string
// 	StderrValue    string
// 	StdStatusVar   string
// 	StdStatusValue int
// 	ResultNames    []string
// 	ResultValues   []interface{}
// }
