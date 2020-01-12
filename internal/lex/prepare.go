package lex

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kassybas/tame/schema"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/loader"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/settings"
	"github.com/sirupsen/logrus"
)

func createContext(globals []tvar.TVariable, sts settings.Settings) (tcontext.Context, error) {
	return tcontext.Context{
		Globals:  globals,
		Settings: sts,
	}, nil
}
func evaluateGlobals(globalDefs map[string]interface{}) ([]tvar.TVariable, error) {
	var vars []tvar.TVariable
	for k, v := range globalDefs {
		if strings.HasSuffix(k, keywords.GlobalDefaultVarSuffix) {
			name := strings.TrimSuffix(k, keywords.GlobalDefaultVarSuffix)
			name = strings.TrimSpace(name)
			sysEnvValue, sysEnvExists := os.LookupEnv(name)
			var value interface{}
			if sysEnvExists {
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

func convertIncludesToRelativePath(path string, includes []schema.IncludeSchema) []schema.IncludeSchema {
	for i := range includes {
		includes[i].Path = fmt.Sprintf("%s%s%s", filepath.Dir(path), string(filepath.Separator), includes[i].Path)
	}
	return includes
}

func PrepareStep(path, targetName string, targetArgs []string) (step.Step, tcontext.Context, error) {
	tf, err := loader.Load(path)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("error loading tamefile: %s\n%s", path, err.Error())
	}
	tf.Includes = convertIncludesToRelativePath(path, tf.Includes)

	root, globalDefs, err := Analyse(tf, targetName, targetArgs)
	if err != nil {
		logrus.Fatal(err)
	}
	globals, err := evaluateGlobals(globalDefs)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("failed to evaluate global variables:\n\t%s", err.Error())
	}
	stgs, err := BuildSettings(tf.Sets)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("failed to evaluate settings:\n\t%s", err.Error())
	}
	ctx, err := createContext(globals, stgs)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("error while creating context:\n\t%s", err.Error())
	}
	return root, ctx, nil
}
