package lex

import (
	"os"

	"github.com/kassybas/mate/schema"

	"github.com/kassybas/mate/internal/helpers"
	"github.com/kassybas/mate/internal/parse"
	"github.com/kassybas/mate/internal/step"
	"github.com/kassybas/mate/internal/tvar"
)

func parseCLITargetArgs(targetArgs []string) ([]tvar.Variable, error) {
	var args []tvar.Variable
	for _, argStr := range targetArgs {
		k, v, err := helpers.GetKeyValueFromEnvString(argStr)
		if err != nil {
			return nil, err
		}
		newArg := tvar.Variable{
			Name:  k,
			Value: v,
		}
		args = append(args, newArg)
	}
	return args, nil
}

func createDependencyGraph(targets map[string]step.Target, targetName string, cliVarArgs []string) (step.Step, error) {
	var root step.Step
	var err error
	root.Name = "root"
	root.CalledTargetName = targetName
	root.Arguments, err = parseCLITargetArgs(cliVarArgs)
	if err != nil {
		return root, err
	}
	root.CalledTarget, err = findCalledTarget(targetName, "cli root", targets)
	if err != nil {
		return root, err
	}

	err = populateSteps(&root.CalledTarget, targets)
	// TODO: continue here

	return root, err
}

// Analyse creates the internal representation
func Analyse(tf schema.Tamefile, targetName string, cliVarArgs []string) (step.Step, map[string]string, error) {

	parsedTargets, err := parse.ParseTeafile(tf)
	if err != nil {
		return step.Step{}, nil, err
	}

	if targetName == "" {
		helpers.PrintTeafileDescription(parsedTargets)
		os.Exit(0)
	}

	// TODO: Load external files if referred

	// build the dependency graph with the called target
	var root step.Step
	root, err = createDependencyGraph(parsedTargets, targetName, cliVarArgs)
	if err != nil {
		return root, nil, err
	}

	return root, tf.Globals, nil
}
