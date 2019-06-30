package lex

import (
	"fmt"
	"os"

	"github.com/kassybas/mate/internal/helpers"
	"github.com/kassybas/mate/internal/loader"
	"github.com/kassybas/mate/internal/parse"
	"github.com/kassybas/mate/types/step"
)

func parseCLITargetArgs(targetArgs []string) ([]step.Argument, error) {
	var args []step.Argument
	for _, argStr := range targetArgs {
		k, v, err := helpers.GetKeyValueFromEnvString(argStr)
		if err != nil {
			return nil, err
		}
		newArg := step.Argument{
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
	root.Arguments, err = parseCLITargetArgs(cliVarArgs)
	if err != nil {
		return root, err
	}
	root.CalledTarget, err = findCalledTarget(targetName, targets)
	if err != nil {
		return root, err
	}

	err = populateSteps(&root.CalledTarget, targets)
	// TODO: continue here
	fmt.Printf("%+v", root)

	return root, err
}

// Analyse creates the internal representation
func Analyse(filePath string, targetName string, cliVarArgs []string) (step.Step, error) {

	tf, err := loader.Load(filePath)
	if err != nil {
		return step.Step{}, err
	}

	parsedTargets, err := parse.ParseTeafile(tf)
	if err != nil {
		return step.Step{}, err
	}

	if targetName == "" {
		helpers.PrintTeafileDescription(parsedTargets)
		os.Exit(0)
	}

	// TODO: Load external files if referred

	// build the dependency graph with the called target
	var head step.Step
	head, err = createDependencyGraph(parsedTargets, targetName, cliVarArgs)
	if err != nil {
		return step.Step{}, err
	}

	return head, nil
}
