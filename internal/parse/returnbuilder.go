package parse

import (
	"fmt"

	"github.com/kassybas/tame/internal/step/returnstep"
)

func buildReturnStep(stepDef map[string]interface{}) (returnstep.ReturnStep, error) {
	var newStep returnstep.ReturnStep
	var err error
	for k, v := range stepDef {
		switch v.(type) {
		case string:
			{
				newStep.Return = []string{v.(string)}
				continue
			}
		case []interface{}:
			{
				newStep.Return, err = ifaceSliceToStringSlice(v.([]interface{}))
				if err != nil {
					return newStep, err
				}
				continue
			}
		default:
			{
				return newStep, fmt.Errorf("unknown type in return step should be string or []string but found: %v (type %T)", v, v)
			}
		}
		return newStep, fmt.Errorf("unknown field in return step: %s", k)
	}
	return newStep, nil
}
