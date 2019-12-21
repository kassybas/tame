package step

import (
	"fmt"

	"github.com/kassybas/shell-exec/exec"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
	"github.com/sirupsen/logrus"
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
	logrus.Fatal("internal error: calling target in shell")
}

func (s *ShellStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) error {
	var err error
	// ignore result if it is not caputered
	ignoreResult := s.Results.StderrVar == "" && s.Results.StdoutVar == "" && s.Results.StdStatusVar == ""
	opts := exec.Options{
		Silent:       s.Opts.Silent,
		ShellPath:    ctx.Settings.UsedShell,
		IgnoreResult: ignoreResult,
		ShieldEnv:    ctx.Settings.ShieldEnv,
	}
	envVars := vt.GetAllEnvVars(ctx.Settings.ShellFieldSeparator)
	prefixedScript := fmt.Sprintf("%s\n%s", ctx.Settings.InitScript, s.Script)
	s.Results.StdoutValue, s.Results.StderrValue, s.Results.StdStatusValue, err = exec.ShellExec(prefixedScript, envVars, opts)
	return err
}
