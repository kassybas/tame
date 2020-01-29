package loadstep

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/kassybas/tame/internal/build/loader"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
	"github.com/kassybas/tame/types/vartype"
	"gopkg.in/yaml.v2"
)

type LoadStep struct {
	basestep.BaseStep
	path          string
	sourceVarName string
	format        string
}

func NewLoadStep(stepDef schema.MergedStepSchema) (*LoadStep, error) {
	var newStep LoadStep
	var err error
	name := fmt.Sprintf("load %s", stepDef.Load.Path)
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.Load, name)
	if stepDef.Load.Path == "" && stepDef.Load.SourceVarName == "" {
		return nil, fmt.Errorf("invalid load step defined: either path or var has to be defined as source of load")
	}
	if stepDef.Load.Path != "" && stepDef.Load.SourceVarName != "" {
		return nil, fmt.Errorf("invalid load step defined: only one source (path or var) can be defined to load")
	}
	newStep.path = stepDef.Load.Path
	newStep.sourceVarName = stepDef.Load.SourceVarName
	newStep.format = stepDef.Load.Format
	return &newStep, err
}

func parseContents(contents string, format string) (interface{}, error) {
	var result interface{}
	switch format {
	case "string", "":
		{
			return contents, nil
		}
	case "yaml":
		{
			err := yaml.Unmarshal([]byte(contents), &result)
			if err != nil {
				return nil, fmt.Errorf("could not encode source variable to yaml in load step:\n\t%s", err.Error())
			}
			return result, nil
		}
	case "json":
		{
			err := json.Unmarshal([]byte(contents), &result)
			if err != nil {
				return nil, fmt.Errorf("could not encode source variable to json in load step:\n\t%s", err.Error())
			}
			return result, nil
		}
	case "toml":
		{
			err := toml.Unmarshal([]byte(contents), &result)
			if err != nil {
				return nil, fmt.Errorf("could not encode source variable to toml in load step:\n\t%s", err.Error())
			}
			return result, nil
		}
	default:
		return nil, fmt.Errorf("unknown encoding in dump step, possible: yaml|json|toml, got: %s", format)
	}

}

func (s *LoadStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	var contents string
	if s.path != "" {
		path, err := vt.ResolveValueToStr(s.path)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("could not resolve expression in path: %s\n\t%s", s.path, err.Error())}
		}
		// load from file
		cbytes, err := loader.ReadFile(path)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("failed to load file in step: %s\n\t%s", s.GetName(), err.Error())}
		}
		contents = string(cbytes)
	} else {
		// load from variable
		v, err := vt.GetVar(s.sourceVarName)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("failed to resolve variable load step: %s\n\t", s.GetName(), err.Error())}
		}
		if v.Type() != vartype.TScalarType {
			return step.StepStatus{Err: fmt.Errorf("only scalar variables can be loaded and parsed in step: %s\n\tgot:%s", s.GetName(), v.Type().Name())}
		}
		contents = v.ToStr()
	}
	resInter, err := parseContents(contents, s.format)
	if err != nil {
		return step.StepStatus{Err: err}
	}
	return step.StepStatus{Results: []interface{}{resInter}}
}
