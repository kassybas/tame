package lex

import (
	"os"

	"github.com/kassybas/tame/schema"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/helpscreen"
	"github.com/kassybas/tame/internal/parse"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tvar"
)

func parseCLITargetArgs(targetArgs []string) ([]tvar.TVariable, error) {
	var args []tvar.TVariable
	for _, argStr := range targetArgs {
		k, v, err := helpers.GetKeyValueFromEnvString(argStr)
		if err != nil {
			return nil, err
		}
		newArg := tvar.NewVariable(k, v)
		args = append(args, newArg)
	}
	return args, nil
}

func createDependencyGraph(targets map[string]step.Target, targetName string, cliVarArgs []string) (step.Step, error) {
	var root step.CallStep
	var err error
	root.Name = targetName
	root.CalledTargetName = targetName
	root.Arguments, err = parseCLITargetArgs(cliVarArgs)
	if err != nil {
		return &root, err
	}
	root.CalledTarget, err = findCalledTarget(targetName, "cli root", targets)
	if err != nil {
		return &root, err
	}

	err = populateSteps(&root.CalledTarget, targets)

	return &root, err
}

// Analyse creates the internal representation
func Analyse(tf schema.Tamefile, targetName string, cliVarArgs []string) (step.Step, map[string]interface{}, error) {

	parsedTargets, err := parse.ParseTeafile(tf)
	if err != nil {
		return nil, nil, err
	}

	if targetName == "" {
		helpscreen.PrintTeafileDescription(parsedTargets)
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
