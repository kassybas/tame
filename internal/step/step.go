package step

import (
	"github.com/kassybas/mate/internal/vartable"

	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
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
	StdoutVar    string
	StdoutValue  string
	StderrVar    string
	StderrValue  string
	StdrcVar     string
	StdrcValue   int
	ResultNames  []string
	ResultValues []interface{}
}
