package orchestrate

import (
	"os"

	"github.com/kassybas/mate/internal/lex"
	"github.com/kassybas/mate/internal/loader"
	"github.com/kassybas/mate/internal/steprunner"
	"github.com/kassybas/mate/types/settings"
	"github.com/kassybas/mate/types/step"
	"github.com/sirupsen/logrus"
)

func EvaluateGlobals(globalDefs map[string]string) ([]step.Variable, error) {
	// TODO
	return nil, nil
}

func CreateContext(globals []step.Variable, sts settings.Settings) (steprunner.Context, error) {
	// TODO
	return steprunner.Context{
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
	_, rc, err := ctx.Run(root.CalledTarget, root.Arguments)
	if err != nil {
		logrus.Fatal("error during execution: ", err.Error())
	}
	// pass through the status code
	os.Exit(rc)
}
