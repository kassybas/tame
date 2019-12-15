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
	GetResult() Result
	GetOpts() opts.ExecutionOpts
	SetOpts(opts.ExecutionOpts)
	RunStep(tcontext.Context, vartable.VarTable) error
	GetCalledTargetName() string
	SetCalledTarget(Target)
}

// TODO: make result a variable interface
type Result struct {
	StdoutVar      string
	StdoutValue    string
	StderrVar      string
	StderrValue    string
	StdStatusVar   string
	StdStatusValue int
	ResultNames    []string
	ResultValues   []interface{}
}
