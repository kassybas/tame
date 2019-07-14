package orchestrate

import (
	"os"

	"github.com/kassybas/mate/internal/vartable"

	"github.com/kassybas/mate/internal/lex"
	"github.com/kassybas/mate/internal/loader"
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/settings"
	"github.com/sirupsen/logrus"
)

func EvaluateGlobals(globalDefs map[string]string) ([]tvar.Variable, error) {
	// TODO
	return nil, nil
}

func CreateContext(globals []tvar.Variable, sts settings.Settings) (tcontext.Context, error) {
	// TODO
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
	os.Exit(root.GetResult().StdrcValue)
}
