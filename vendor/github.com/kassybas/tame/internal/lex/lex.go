package lex

import (
	"github.com/kassybas/tame/internal/dependency"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/loader"
	"github.com/kassybas/tame/internal/parser"
	"github.com/kassybas/tame/internal/target"
	"os"
)

func parseCLITargetArgs(targetArgs []string)(map[string]string, error){
	argsMap := make(map[string]string)
	for _, argStr := range targetArgs {
		k, v, err := helpers.GetKeyValueFromEnvString(argStr)
		if err!=nil{
			return nil, err
		}
		argsMap[k] = v
	}
	return argsMap, nil
}

func createDependencyGraph(targets map[string]target.Target, trg target.Target, targetArgs []string) (dependency.Dependency, error) {
	var depConf target.DepConfig
	var err error
	depConf.ArgValues, err = parseCLITargetArgs(targetArgs)
	if err != nil {
		return dependency.Dependency{}, err
	}

	head, err := createDependencyFromTargets(trg, depConf, targets)
	if err != nil {
		return dependency.Dependency{}, err
	}
	return head, nil
}

// Analyse creates the internal representation
func Analyse(filePath string, targetName string, targetArgs []string) (dependency.Dependency, error) {

	tf, err := loader.Load(filePath)
	if err != nil {
		return dependency.Dependency{}, err
	}

	parsedTargets, err := parser.ParseTeafile(tf)
	if err != nil {
		return dependency.Dependency{}, err
	}

	if targetName == "" {
		helpers.PrintTeafileDescription(parsedTargets)
		os.Exit(0)
	}

	// Find target called from CLI
	trg, err := findCalledTarget(targetName, parsedTargets)
	if err != nil {
		return dependency.Dependency{}, err
	}
	// TODO: Load external files if referred

	// build the dependency graph with the called target
	var head dependency.Dependency
	head, err = createDependencyGraph(parsedTargets, trg, targetArgs)
	if err != nil {
		return dependency.Dependency{}, err
	}

	return head, nil
}
