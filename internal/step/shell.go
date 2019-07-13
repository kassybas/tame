package step

import (
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
	"github.com/kassybas/shell-exec/exec"
)

type Shell struct {
	Name    string
	Opts    opts.ExecutionOpts
	Results Result
	Script  string
}

func (s *Shell) GetName() string {
	return ""
}

func (s *Shell) Kind() steptype.Steptype {
	return steptype.Shell
}

func (s *Shell) SetOpts(o opts.ExecutionOpts) {
	s.Opts = o
}

func (s *Shell) GetResult() Result {
	return s.Results
}
func (s *Shell) RunStep(ctx tcontext.Context, vars map[string]tvar.Variable) ([]string, Result, error) {
	var err error
	// ignore result if neither stdout variable and stderr variable is defined
	ignoreResult := s.Results.StderrVar == "" && s.Results.StdoutVar == ""
	opts := exec.Options{
		Silent:       s.Opts.Silent,
		ShellPath:    ctx.Settings.UsedShell,
		IgnoreResult: ignoreResult,
	}
	envVars := FormatEnvVars(vars)
	prefixedScript := ctx.Settings.InitScript + "\n" + s.Script
	s.Results.StdoutValue, s.Results.StderrValue, s.Results.StdrcValue, err = exec.ShellExec(prefixedScript, envVars, opts)
	return nil, s.Results, err
}
