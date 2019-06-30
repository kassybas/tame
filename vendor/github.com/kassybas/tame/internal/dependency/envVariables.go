package dependency

import (
	"fmt"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"strings"
)

func (td *Dependency) CreateEnvVarsFromDeps() ([]string, error) {
	var envVars []string
	// TODO check if dep (regardless of trg) executed
	fullName := helpers.FlattenEnvVarNameLocal(td.Name)
	if td.ExecOpts.SaveOut {
		envVars = append(envVars, keywords.PrefixOut+fullName+"="+td.stdOut)
	}
	if td.ExecOpts.SaveErr {
		envVars = append(envVars, keywords.PrefixErr+fullName+"="+td.stdErr)
	}
	if td.ExecOpts.SaveRc {
		envVars = append(envVars, keywords.PrefixRc+fullName+"="+td.stdRc)
	}
	// TODO: check if all given args were consumed

	return envVars, nil
}

func (td *Dependency) CreateEnvVarsFromArgs() ([]string, error) {
	var envVars []string
	for _, param := range td.Target.Params {
		// Check given args
		argVal, err := td.GetArgValue(param.Name)
		if err != nil {
			return nil, err
		}
		paramEnvVarName := helpers.FlattenEnvVarNameLocal(param.Name)
		if err != nil {
			return nil, err
		}
		envVars = append(envVars, paramEnvVarName+"="+argVal)
		continue
	}
	return envVars, nil
}



func (td *Dependency) CreateEnvVarsFromGlobals() ([]string, error) {
	var envVars []string

	for _, g := range *td.Target.GlobalVars {
		var s string
		s = g.EnvVarName + "=" + g.Value
		envVars = append(envVars, s)

	}
	return envVars, nil
}

func (td *Dependency) GetReferedValue(refStr string) (string, error) {
	if td.Caller == nil {
		return "", fmt.Errorf("cannot use variable reference on top level target: '%s'", refStr)
	}

	trimmedRefStr := refStr[1:]
	v, err := td.Caller.GetArgValue(trimmedRefStr)
	if err != nil {
		return "", fmt.Errorf("refered variable not found: '%s'", refStr)
	}
	return v, nil
}

func (td *Dependency) GetArgValue(argName string) (string, error) {
	var err error
	val, exists := td.ArgValues[argName]

	if !exists {
		val, err = td.Target.GetDefaultValue(argName)
		if err != nil {
			return "", err
		}
	}
	return val, nil
}

func (td *Dependency) GetAllEnvVars() ([]string, error) {
	var allEnvVars []string
	var envVars []string
	var err error

	// 1. GlobalVars
	envVars, err = td.CreateEnvVarsFromGlobals()
	if err != nil {
		return allEnvVars, err
	}
	allEnvVars = append(allEnvVars, envVars...)

	// 2. Deps
	for _, d := range td.Deps {
		envVars, err := d.CreateEnvVarsFromDeps()
		if err != nil {
			return allEnvVars, err
		}
		allEnvVars = append(allEnvVars, envVars...)
	}

	// 3. Args
	envVars, err = td.CreateEnvVarsFromArgs()
	if err != nil {
		return allEnvVars, err
	}
	allEnvVars = append(allEnvVars, envVars...)

	return allEnvVars, nil
}

func (td *Dependency) EvaluateReferences() error {
	var err error
	for k, val := range td.ArgValues {
		if strings.HasPrefix(val, keywords.PrefixReference) {
			td.ArgValues[k], err = td.GetReferedValue(val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
