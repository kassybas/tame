package shellstep

import (
	"fmt"

	"github.com/kassybas/shell-exec/exec"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
)

func (s *ShellStep) shouldIgnoreResults() bool {
	if len(s.BaseStep.ResultNames()) == 0 {
		return true
	}
	if len(s.BaseStep.ResultNames()) == 1 {
		if s.BaseStep.ResultNames()[0] != "" {
			return false
		}
		return true
	}
	if s.BaseStep.ResultNames()[0] == "" && s.BaseStep.ResultNames()[1] == "" {
		return true
	}
	return false
}

func (s *ShellStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	var err error
	// ignore result if it is not caputered
	// TODO: fix regression
	opts := exec.Options{
		Silent:       s.BaseStep.GetOpts().Silent,
		ShellPath:    ctx.Settings.UsedShell,
		IgnoreResult: s.shouldIgnoreResults(),
		ShieldEnv:    ctx.Settings.ShieldEnv,
	}
	envVars := vt.GetAllEnvVars(ctx.Settings.ShellFieldSeparator)
	prefixedScript := fmt.Sprintf("%s\n%s", ctx.Settings.InitScript, s.script)
	stdoutValue, stderrValue, stdStatusValue, err := exec.ShellExec(prefixedScript, envVars, opts)
	return step.StepStatus{
		Results:   []interface{}{stdoutValue, stderrValue, stdStatusValue},
		Stdstatus: stdStatusValue,
		Err:       err,
	}
}
