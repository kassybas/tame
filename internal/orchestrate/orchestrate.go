package orchestrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/kassybas/tame/internal/keywords"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/kassybas/tame/internal/lex"
	"github.com/kassybas/tame/internal/loader"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/settings"
	"github.com/sirupsen/logrus"
)

func EvaluateGlobals(globalDefs map[string]interface{}) ([]tvar.TVariable, error) {
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

func CreateContext(globals []tvar.TVariable, sts settings.Settings) (tcontext.Context, error) {
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
		logrus.Fatal("failed to evaluate global variables", err.Error())
	}
	stgs, err := lex.BuildSettings(tf.Sets)
	if err != nil {
		logrus.Fatal("failed to evaluate settings", err.Error())
	}
	ctx, err := CreateContext(globals, stgs)
	if err != nil {
		logrus.Fatal("error while creating context", err.Error())
	}

	// TODO: put cli args in here
	status := root.RunStep(ctx, vartable.NewVarTable())
	if status.Err != nil {
		logrus.Fatal("error:\n\t", status.Err.Error())
	}
	// pass through the status code
	if status.Stdstatus != 0 {
		fmt.Fprintf(os.Stdout, "tame: *** [%s] Error %d\n", targetName, status.Stdstatus)
	}

	os.Exit(status.Stdstatus)
}
