package step

import (
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/vartable"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
	"github.com/kassybas/shell-exec/exec"
)

type ShellStep struct {
	Name    string
	Opts    opts.ExecutionOpts
	Results Result
	Script  string
}

func (s *ShellStep) GetOpts() opts.ExecutionOpts {
	return s.Opts
}

func (s *ShellStep) GetName() string {
	return ""
}

func (s *ShellStep) Kind() steptype.Steptype {
	return steptype.Shell
}

func (s *ShellStep) SetOpts(o opts.ExecutionOpts) {
	s.Opts = o
}

func (s *ShellStep) GetCalledTargetName() string {
	return "shell"
}

func (s *ShellStep) GetResult() Result {
	return s.Results
}

func (s *ShellStep) SetCalledTarget(t Target) {
	panic("calling target in shell")
}

func (s *ShellStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) error {
	var err error
	// ignore result if neither stdout variable and stderr variable is defined
	ignoreResult := s.Results.StderrVar == "" && s.Results.StdoutVar == ""
	opts := exec.Options{
		Silent:       s.Opts.Silent,
		ShellPath:    ctx.Settings.UsedShell,
		IgnoreResult: ignoreResult,
		ShieldEnv:    ctx.Settings.ShieldEnv,
	}
	envVars := vt.GetAllEnvVars()
	prefixedScript := ctx.Settings.InitScript + "\n" + s.Script
	s.Results.StdoutValue, s.Results.StderrValue, s.Results.StdrcValue, err = exec.ShellExec(prefixedScript, envVars, opts)

	return err
}
