package orchestrate

import (
	"os"
	"strings"

	"github.com/kassybas/mate/internal/keywords"

	"github.com/kassybas/mate/internal/vartable"

	"github.com/kassybas/mate/internal/lex"
	"github.com/kassybas/mate/internal/loader"
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/settings"
	"github.com/sirupsen/logrus"
)

func EvaluateGlobals(globalDefs map[string]interface{}) ([]tvar.VariableI, error) {
	var vars []tvar.VariableI
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
			vars = append(vars, tvar.CreateVariable(name, value))
			continue
		}
		vars = append(vars, tvar.CreateVariable(k, v))
	}
	return vars, nil
}

func CreateContext(globals []tvar.VariableI, sts settings.Settings) (tcontext.Context, error) {
	return tcontext.Context{
		Globals:  globals,
		Settings: sts,
	}, nil
}

func Make(path, targetName string, targetArgs []string) {
	tf, err := loader.Load(path)
	if err != nil {
		logrus.Fatalf("error loading tamefile: %s\n%s", path, err.Error())
	}
	root, globalDefs, err := lex.Analyse(tf, targetName, targetArgs)
	if err != nil {
		logrus.Fatal(err)
	}
	globals, err := EvaluateGlobals(globalDefs)
	if err != nil {
		logrus.Fatal("error while evaluating global variables", err.Error())
	}
	stgs, err := lex.BuildSettings(tf.Sets)
	if err != nil {
		logrus.Fatal("error while evaluating settings", err.Error())
	}
	ctx, err := CreateContext(globals, stgs)
	if err != nil {
		logrus.Fatal("error while creating context", err.Error())
	}

	// TODO: put cli args in here
	err = root.RunStep(ctx, vartable.NewVarTable())
	if err != nil {
		logrus.Fatal("error during execution: ", err.Error())
	}
	// pass through the status code
	os.Exit(root.GetResult().StdStatusValue)
}
