package dumpstep

import (
	"fmt"
	"io"
	"os"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type DumpStep struct {
	basestep.BaseStep
	sourceValue string
	path        string
	format      string
}

func NewDumpStep(stepDef schema.MergedStepSchema) (*DumpStep, error) {
	var newStep DumpStep
	var err error
	name := fmt.Sprintf("dump %s -> %s [format: %s]", stepDef.Dump.SourceValue, stepDef.Dump.Path, stepDef.Dump.Format)
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Wait, name)
	newStep.sourceValue = stepDef.Dump.SourceValue
	newStep.format = stepDef.Dump.Format
	newStep.path = stepDef.Dump.Path
	return &newStep, err
}

func writeToFile(path string, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func (s *DumpStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	var err error
	sourceVal, err := vt.ResolveValue(s.sourceValue)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("source value cannot be resolved dump step: %s\n\t%s", s.GetName(), err.Error())}
	}
	dumpedValue, err := helpers.GetFormattedValue(sourceVal, s.format)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("failed to encode in dump step: %s\n\t%s", s.GetName(), err.Error())}
	}
	if s.path != "" {
		path, err := vt.ResolveValueToStr(s.path)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("could not resolve expression in path: %s\n\t%s", s.path, err.Error())}
		}
		err = writeToFile(path, dumpedValue)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("failed to save file in dump step: %s\n\t%s", s.GetName(), err.Error())}
		}
	}
	if !s.GetOpts().Silent {
		fmt.Print(dumpedValue)
	}
	return step.StepStatus{Results: []interface{}{dumpedValue}}
}
