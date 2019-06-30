package parser

import (
	"fmt"
	"github.com/kassybas/tame/internal/target"
)

func buildArgumentDefinition(argValue interface{}) (target.ParamConfig, error) {
	var arg target.ParamConfig

	// Multiple types possible in schema. Eg:
	// target-name:
	//   args:
	//   - arg1                   // <- no default value:  string
	//   - arg2: default-value1   // <- default value set: map[string]string
	switch argValue.(type) {
	case string:
		{
			arg.Name = argValue.(string)
			return arg, nil
		}
	case map[interface{}]interface{}:
		{
			for k, v := range argValue.(map[interface{}]interface{}) {
				arg.Name = k.(string)
				arg.HasDefault = true
				arg.DefaultValue = v.(string)
			}
			return arg, nil
		}
	}
	return arg, fmt.Errorf("incorrect type for argument definition: %s", argValue)
}