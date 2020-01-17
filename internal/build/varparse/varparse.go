package varparse

import (
	"os"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/tvar"
)

func EvaluateGlobals(globalDefs map[string]interface{}) ([]tvar.TVariable, error) {
	var vars []tvar.TVariable
	for k, v := range globalDefs {
		if strings.HasSuffix(k, keywords.GlobalDefaultVarSuffix) {
			name := strings.TrimSuffix(k, keywords.GlobalDefaultVarSuffix)
			name = strings.TrimSpace(name)
			var value interface{}
			if sysEnvValue, sysEnvExists := os.LookupEnv(name); sysEnvExists {
				value = sysEnvValue
			} else {
				value = v
			}
			vars = append(vars, tvar.NewVariable(name, value))
			continue
		}
		vars = append(vars, tvar.NewVariable(k, v))
	}
	return vars, nil
}

func ParseGlobals(dynamicKeys map[string]interface{}) (map[string]interface{}, error) {
	globals := map[string]interface{}{}
	for k, v := range dynamicKeys {
		if !strings.HasPrefix(k, keywords.PrefixReference) {
			continue
		}
		globals[k] = v
		delete(dynamicKeys, k)
	}
	return globals, nil
}
