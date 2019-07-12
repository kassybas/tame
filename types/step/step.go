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
	Arguments        []Variable
	Opts             opts.ExecutionOpts
	HasResult        bool
	Results          Result
	CalledTargetName string
	CalledTarget     Target
	Script           string
}

type Result struct {
	StdoutVar    string
	StdoutValue  string
	StderrVar    string
	StderrValue  string
	StdrcVar     string
	StdrcValue   int
	ResultVars   []string
	ResultValues []string
}
