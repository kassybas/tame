package stepparse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
	"github.com/mitchellh/mapstructure"
)

func validateVariableName(name string) error {
	if !strings.HasPrefix(name, keywords.PrefixReference) {
		return fmt.Errorf("variables and arguments must start with '$' symbol: %s (correct: %s%s)", name, keywords.PrefixReference, name)
	}
	return nil
}
func parseCallStepArgs(argDefs interface{}) (map[string]interface{}, error) {
	argMap, ok := argDefs.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("call step must have map as value, got: %T", argDefs)
	}
	args := make(map[string]interface{}, len(argMap))
	for argKey, argValue := range argMap {
		argName, ok := argKey.(string)
		if !ok {
			return nil, fmt.Errorf("non-string argument variable name: %v (type %T)", argKey, argKey)
		}
		if err := validateVariableName(argName); err != nil {
			return nil, err
		}
		args[argName] = argValue
	}
	return args, nil
}

func parseCalledTargetName(k string) (string, error) {
	fields := strings.Fields(k)
	if len(fields) == 0 {
		return "", fmt.Errorf("'%s': no called target name found", k)
	}
	if len(fields) > 1 {
		return "", fmt.Errorf("'%s': called target name contains whitespaces", k)
	}
	return fields[0], nil
}

func loadCallStepSchema(raw map[string]interface{}, dynamicKey string, result *schema.MergedStepSchema) error {
	ct, err := parseCalledTargetName(dynamicKey)
	if err != nil {
		return err
	}
	result.CalledTargetName = ct
	result.CallArgumentsPassed, err = parseCallStepArgs(raw[dynamicKey])
	if err != nil {
		return fmt.Errorf("error while parsing arguments of step '%s' args: '%v'\n\t%s", result.CalledTargetName, result.CallArgumentsPassed, err.Error())
	}
	return err
}

func loadVarStepSchema(raw map[string]interface{}, dynamicKey string, result *schema.MergedStepSchema) error {
	err := validateVariableName(dynamicKey)
	if err != nil {
		return fmt.Errorf("error while parsing step '%s'\n\t%s", dynamicKey, err.Error())
	}
	result.VarName = dynamicKey
	result.ResultContainers = &[]string{dynamicKey}
	result.VarValue = raw[dynamicKey]
	return err
}

func loadSubSteps(rawSteps interface{}) ([]schema.MergedStepSchema, error) {
	rawStepsList, ok := rawSteps.([]interface{})
	if !ok {
		return nil, fmt.Errorf("substeps must be lists, got: %T", rawSteps)
	}
	steps := make([]schema.MergedStepSchema, len(rawStepsList))
	for i := range rawStepsList {
		stepMap, err := helpers.ConvertInterToMapStrInter(rawStepsList[i])
		steps[i], err = ParseStepSchema(stepMap)
		if err != nil {
			return nil, err
		}
	}
	return steps, nil
}

func loadIfStepSchema(raw map[string]interface{}, result *schema.MergedStepSchema) error {
	var err error
	if len(raw) > 2 {
		return fmt.Errorf("unknown key(s) in if-else step: %v", raw)
	}
	for k, v := range raw {
		if strings.HasPrefix(k, keywords.IfStepPrefix) {
			condition := strings.TrimPrefix(k, keywords.IfStepPrefix)
			result.IfCondition = condition
			if result.IfSteps, err = loadSubSteps(v); err != nil {
				return err
			}
		} else if k == keywords.IfStepElse {
			if result.ElseSteps, err = loadSubSteps(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadDynamicKey(raw map[string]interface{}, dynamicKey string, result *schema.MergedStepSchema) error {
	// starting with $ -> it's a var
	if strings.HasPrefix(dynamicKey, keywords.PrefixReference) && len(dynamicKey) > 1 {
		return loadVarStepSchema(raw, dynamicKey, result)
	}
	if strings.HasPrefix(dynamicKey, keywords.IfStepPrefix) || dynamicKey == keywords.IfStepElse {
		return loadIfStepSchema(raw, result)
	}
	// not starting with $ -> it's a call
	return loadCallStepSchema(raw, dynamicKey, result)
}

func loadMergedStepSchema(raw interface{}) (schema.MergedStepSchema, error) {
	var result schema.MergedStepSchema
	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("step must be a map, got %T", raw)
	}
	var md mapstructure.Metadata
	err := mapstructure.WeakDecodeMetadata(raw, &result, &md)
	if err != nil {
		return result, err
	}
	// 2: if and else
	if len(md.Unused) > 2 {
		return result, fmt.Errorf("multiple unknown keys found in step, expect: $varname, targetname, if, else; got: %v", md.Unused)
	}
	if len(md.Unused) > 0 {
		err = loadDynamicKey(rawMap, md.Unused[0], &result)
		if err != nil {
			return result, err
		}
	}
	return result, err
}

func setStepType(oldStepType, newStepType steptype.Steptype) (steptype.Steptype, error) {
	if oldStepType != steptype.Unset {
		return 0, fmt.Errorf("multiple step types found: %s and %s", oldStepType.ToStr(), newStepType.ToStr())
	}
	return newStepType, nil
}

func ParseStepSchema(raw interface{}) (schema.MergedStepSchema, error) {
	var err error
	mergedSchema, err := loadMergedStepSchema(raw)
	if err != nil {
		return mergedSchema, err
	}
	// Determine step type
	if mergedSchema.Print != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Print)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.Load != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Load)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.Dump != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Dump)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.ForLoop != nil {
		if mergedSchema.ForSteps, err = loadSubSteps(mergedSchema.ForRawSteps); err != nil {
			return mergedSchema, fmt.Errorf("failed to parse for-do block steps: %s", err.Error())
		}
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.For)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.IfCondition != "" {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.If)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.Wait != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Wait)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.Expr != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Expr)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.Return != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Return)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.VarName != "" {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Var)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.CalledTargetName != "" {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Call)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.Script != nil {
		mergedSchema.StepType, err = setStepType(mergedSchema.StepType, steptype.Shell)
		if err != nil {
			return mergedSchema, err
		}
	}
	if mergedSchema.StepType == steptype.Unset {
		return mergedSchema, fmt.Errorf("could not determine step type: %+v", mergedSchema)
	}
	return mergedSchema, err
}
