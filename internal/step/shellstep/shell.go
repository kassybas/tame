package shellstep

import (
	"fmt"
	"strings"

	"github.com/kassybas/shell-exec/exec"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type ShellStep struct {
	basestep.BaseStep
	script string
}

func NewShellStep(stepDef schema.MergedStepSchema) (*ShellStep, error) {
	var err error
	var newStep ShellStep
	if stepDef.Script == nil {
		return &newStep, fmt.Errorf("missing called script in shell step")
	}
	newStep.script = strings.Join(*stepDef.Script, "\n")
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Shell, "shell")
	return &newStep, err
}

func (s *ShellStep) shouldIgnoreResults() bool {
	if len(s.BaseStep.ResultNames()) == 0 {
		return true
	}
	if len(s.BaseStep.ResultNames()) == 1 {
		// igrnore if empty
		return s.BaseStep.ResultNames()[0] == ""
	}
	return s.BaseStep.ResultNames()[0] == "" && s.BaseStep.ResultNames()[1] == ""
}

func (s *ShellStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	var err error
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
