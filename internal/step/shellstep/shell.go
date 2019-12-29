package shellstep

import (
	"fmt"

	"github.com/kassybas/shell-exec/exec"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
	"github.com/sirupsen/logrus"
)

type ShellStep struct {
	Name        string
	Opts        opts.ExecutionOpts
	Script      string
	Results     []string
	IteratorVar string
	IterableVar string
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

func (s *ShellStep) ResultNames() []string {
	return s.Results
}

func (s *ShellStep) SetCalledTarget(t interface{}) {
	logrus.Fatal("internal error: calling target in shell")
}

func (s *ShellStep) shouldIgnoreResults() bool {
	if len(s.Results) == 0 {
		return true
	}
	if len(s.Results) == 1 {
		if s.Results[0] != "" {
			return false
		}
		return true
	}
	if s.Results[0] == "" && s.Results[1] == "" {
		return true
	}
	return false
}

func (s *ShellStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	var err error
	// ignore result if it is not caputered
	// TODO: fix regression
	opts := exec.Options{
		Silent:       s.Opts.Silent,
		ShellPath:    ctx.Settings.UsedShell,
		IgnoreResult: s.shouldIgnoreResults(),
		ShieldEnv:    ctx.Settings.ShieldEnv,
	}
	envVars := vt.GetAllEnvVars(ctx.Settings.ShellFieldSeparator)
	prefixedScript := fmt.Sprintf("%s\n%s", ctx.Settings.InitScript, s.Script)
	stdoutValue, stderrValue, stdStatusValue, err := exec.ShellExec(prefixedScript, envVars, opts)
	return step.StepStatus{
		Results:   []interface{}{stdoutValue, stderrValue, stdStatusValue},
		Stdstatus: stdStatusValue,
		Err:       err,
	}
}

func (s *ShellStep) GetIteratorVar() string {
	return s.IteratorVar
}

func (s *ShellStep) GetIterableVar() string {
	return s.IterableVar
}
