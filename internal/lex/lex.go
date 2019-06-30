// package lex

// import (
// 	"os"

// 	"github.com/kassybas/mate/internal/dependency"
// 	"github.com/kassybas/mate/internal/helpers"
// 	"github.com/kassybas/mate/internal/loader"
// 	"github.com/kassybas/mate/internal/parser"
// )

// // Analyse creates the internal representation
// func Analyse(filePath string, targetName string, targetArgs []string) (dependency.Dependency, error) {

// 	tf, err := loader.Load(filePath)
// 	if err != nil {
// 		return dependency.Dependency{}, err
// 	}

// 	parsedTargets, err := parser.ParseTeafile(tf)
// 	if err != nil {
// 		return dependency.Dependency{}, err
// 	}

// 	if targetName == "" {
// 		helpers.PrintTeafileDescription(parsedTargets)
// 		os.Exit(0)
// 	}

// 	// Find target called from CLI
// 	trg, err := findCalledTarget(targetName, parsedTargets)
// 	if err != nil {
// 		return dependency.Dependency{}, err
// 	}
// 	// TODO: Load external files if referred

// 	// build the dependency graph with the called target
// 	var head dependency.Dependency
// 	head, err = createDependencyGraph(parsedTargets, trg, targetArgs)
// 	if err != nil {
// 		return dependency.Dependency{}, err
// 	}

// 	return head, nil
// }
