package dumpstep

import (
	"encoding/json"
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
	"gopkg.in/yaml.v2"
)

type DumpStep struct {
	basestep.BaseStep
	sourceVarName string
	path          string
	format        string
}

func NewDumpStep(stepDef schema.MergedStepSchema) (*DumpStep, error) {
	var newStep DumpStep
	var err error
	name := fmt.Sprintf("dump %s -> %s [format: %s]", stepDef.Dump.SourceVarName, stepDef.Dump.Path, stepDef.Dump.Format)
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Wait, name)
	newStep.sourceVarName = stepDef.Dump.SourceVarName
	newStep.format = stepDef.Dump.Format
	newStep.path = stepDef.Dump.Path
	return &newStep, err
}

func getFormattedValue(v interface{}, format string) (string, error) {
	var dumpedValue []byte
	var err error
	switch format {
	case "yaml", "":
		{
			dumpedValue, err = yaml.Marshal(&v)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to yaml in dump step: %s", err.Error())
			}
		}
	case "json":
		{
			// json unmarshal does not support map[interface{}] interface
			// so we need to convert it to map[string]interface{}
			strMapValue, err := helpers.DeepConvertInterToMapStrInter(v)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to json in dump step: %s", err.Error())
			}
			dumpedValue, err = json.Marshal(&strMapValue)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to json in dump step: %s", err.Error())
			}
		}
	default:
		return "", fmt.Errorf("unknown encoding in dump step, possible: yaml|json, got: %s", format)
	}
	return string(dumpedValue), nil

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
	source, err := vt.GetVar(s.sourceVarName)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("source variable does not exist in dump step: %s\n\t%s", s.GetName(), err.Error())}
	}
	dumpedValue, err := getFormattedValue(source.Value(), s.format)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("failed to encode in dump step: %s\n\t%s", s.GetName(), err.Error())}
	}
	if s.path != "" {
		err = writeToFile(s.path, dumpedValue)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("failed to save file in dump step: %s\n\t%s", s.GetName(), err.Error())}
		}
	}
	return step.StepStatus{Results: []interface{}{dumpedValue}}
}
