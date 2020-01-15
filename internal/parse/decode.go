package parse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
	"github.com/mitchellh/mapstructure"
)

func loadCallStepSchema(raw map[string]interface{}, dynamicKey string, result *schema.MergedStepSchema) error {
	ct, err := parseCalledTargetName(dynamicKey)
	if err != nil {
		return err
	}
	result.CalledTargetName = &ct
	result.CallArgumentsPassed, err = parseCallStepArgs(raw[dynamicKey])
	if err != nil {
		return fmt.Errorf("error while parsing arguments of call step [%s]\n\t%s", result.CallArgumentsPassed, err.Error())
	}
	return err
}

func loadVarStepSchema(raw map[string]interface{}, dynamicKey string, result *schema.MergedStepSchema) error {
	err := validateVariableName(dynamicKey)
	if err != nil {
		return fmt.Errorf("error while parsing step '%s'\n\t%s", dynamicKey, err.Error())
	}
	result.VarName = &dynamicKey
	result.ResultContainers = &[]string{dynamicKey}
	result.VarValue = raw[dynamicKey]
	return err
}

func loadDynamicKey(raw map[string]interface{}, dynamicKey string, result *schema.MergedStepSchema) error {
	// starting with $ -> it's a var
	if strings.HasPrefix(dynamicKey, keywords.PrefixReference) && len(dynamicKey) > 1 {
		return loadVarStepSchema(raw, dynamicKey, result)
	}
	// not starting with $ -> it's a call
	return loadCallStepSchema(raw, dynamicKey, result)
}

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
		err = loadDynamicKey(raw, md.Unused[0], &result)
		if err != nil {
			return result, err
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

	if mergedSchema.Wait != nil {
		stepType, err = setStepType(stepType, steptype.Wait)
		if err != nil {
			return mergedSchema, 0, err
		}
	}
	if mergedSchema.Expr != nil {
		stepType, err = setStepType(stepType, steptype.Expr)
		if err != nil {
			return mergedSchema, 0, err
		}
	}
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
	if stepType == steptype.Unset {
		return mergedSchema, stepType, fmt.Errorf("could not determine step type: %+v", mergedSchema)
	}
	return mergedSchema, stepType, err
}
