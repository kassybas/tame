package parse

import (
	"fmt"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
	"github.com/mitchellh/mapstructure"
	"strings"
)

func loadMergedStepSchema(raw map[string]interface{}) (schema.MergedStepSchema, error) {
	var result schema.MergedStepSchema
	var md mapstructure.Metadata
	err := mapstructure.WeakDecodeMetadata(raw, &result, &md)
	if err != nil {
		return result, err
	}
	if len(md.Unused) > 1 {
		return result, fmt.Errorf("multiple dyanmic keys found in step, only one allowed (var ... or call ...), got: %v", md.Unused)
	}
	if len(md.Unused) == 1 {
		dynamicKey := md.Unused[0]
		if strings.HasPrefix(dynamicKey, keywords.StepCall) {
			ct, err := parseCalledTargetName(dynamicKey)
			if err != nil {
				return result, err
			}
			result.CalledTargetName = &ct

			result.CallArgumentsPassed, err = parseCallStepArgs(raw[dynamicKey])
			if err != nil {
				return result, err
			}
		} else if strings.HasPrefix(dynamicKey, keywords.StepVar) {
			result.VarName = &dynamicKey
			result.VarValue = raw[dynamicKey]
		} else {
			return result, fmt.Errorf("unknown key in step: %s", dynamicKey)
		}
	}
	return result, err
}

func setStepType(oldStepType, newStepType steptype.Steptype) (steptype.Steptype, error) {
	if oldStepType != steptype.Unset {
		return 0, fmt.Errorf("multiple step types definied: %s and %s", oldStepType.ToStr(), newStepType.ToStr())
	}
	return newStepType, nil
}

func ParseStepSchema(raw map[string]interface{}) (schema.MergedStepSchema, steptype.Steptype, error) {
	var err error
	mergedSchema, err := loadMergedStepSchema(raw)
	if err != nil {
		return mergedSchema, 0, err
	}
	// Determine step type
	var stepType steptype.Steptype
	if mergedSchema.Return != nil {
		stepType, err = setStepType(stepType, steptype.Return)
		if err != nil {
			return mergedSchema, 0, err
		}
	}
	if mergedSchema.VarName != nil {
		stepType, err = setStepType(stepType, steptype.Var)
		if err != nil {
			return mergedSchema, 0, err
		}
	}
	if mergedSchema.CalledTargetName != nil {
		stepType, err = setStepType(stepType, steptype.Call)
		if err != nil {
			return mergedSchema, 0, err
		}
	}
	if mergedSchema.Script != nil {
		stepType, err = setStepType(stepType, steptype.Shell)
		if err != nil {
			return mergedSchema, 0, err
		}
	}
	return mergedSchema, stepType, err
}
