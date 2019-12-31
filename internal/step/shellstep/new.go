package shellstep

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/step/basestep"
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
