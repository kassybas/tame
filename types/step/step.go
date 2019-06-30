package step

import (
	"github.com/kassybas/mate/types/opts"
)

type Steptype int

const (
	Unset Steptype = iota
	Call
	Exec
)

type Step struct {
	Name             string
	Kind             Steptype
	Arguments        []Argument
	Opts             opts.ExecutionOpts
	HasResult        bool
	ResultVars       Result
	CalledTargetName string
	CalledTarget     Target
	Script           string
}

type Argument struct {
	Name  string
	Value string
}

type Result struct {
	StdoutVar string
	StderrVar string
	StdrcVar  string
}
