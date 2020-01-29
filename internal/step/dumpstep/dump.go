package dumpstep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
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

func getFormattedValue(v interface{}, format string) (string, error) {
	var dumpedValue string
	switch format {
	case "yaml", "":
		{
			dumpedValueBytes, err := yaml.Marshal(&v)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to yaml in dump step: %s", err.Error())
			}
			dumpedValue = string(dumpedValueBytes)
		}
	case "json":
		{
			// json does not support map[interface{}]interface{}
			// so we need to convert it to map[string]interface{}
			strMapValue, err := helpers.DeepConvertInterToMapStrInter(v)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to json in dump step: %s", err.Error())
			}
			dumpedValueBytes, err := json.Marshal(&strMapValue)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to json in dump step: %s", err.Error())
			}
			dumpedValue = string(dumpedValueBytes)
		}
	case "toml":
		{
			// toml does not support map[interface{}] interface{}
			// so we need to convert it to map[string]interface{}
			strMapValue, err := helpers.DeepConvertInterToMapStrInter(v)
			if err != nil {
				return "", fmt.Errorf("could not encode source variable to toml in dump step: %s", err.Error())
			}
			buf := new(bytes.Buffer)
			if err := toml.NewEncoder(buf).Encode(strMapValue); err != nil {
				return "", fmt.Errorf("could not encode source variable to toml in dump step: %s", err.Error())
			}
			dumpedValue = buf.String()
		}
	default:
		return "", fmt.Errorf("unknown encoding in dump step, possible: yaml|json|toml, got: %s", format)
	}
	return dumpedValue, nil

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
		return step.StepStatus{Err: fmt.Errorf("source variable cannot be resolved dump step: %s\n\t%s", s.GetName(), err.Error())}
	}
	dumpedValue, err := getFormattedValue(sourceVal, s.format)
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
		fmt.Println(dumpedValue)
	}
	return step.StepStatus{Results: []interface{}{dumpedValue}}
}
