package orchestrate

import (
	"github.com/kassybas/mate/internal/lex"
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
	root, globalDefs, err := lex.Analyse(path, targetName, targetArgs)
	if err != nil {
		logrus.Fatal(err)
	}
	globals, err := EvaluateGlobals(globalDefs)
	ctx, err := CreateContext(globals, settings.Settings{})
	if err != nil {
		logrus.Fatal("Execution error:", err.Error())
	}

	err = ctx.Exec(root.CalledTarget, root.Arguments)
	if err != nil {
		logrus.Fatal("Execution error:", err.Error())
	}
}
