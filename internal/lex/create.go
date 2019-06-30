// package lex

// import (
// 	"github.com/kassybas/mate/internal/dependency"
// 	"github.com/kassybas/mate/internal/target"
// )

// func parseCLITargetArgs(targetArgs []string) (map[string]string, error) {
// 	argsMap := make(map[string]string)
// 	for _, argStr := range targetArgs {
// 		k, v, err := helpers.GetKeyValueFromEnvString(argStr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		argsMap[k] = v
// 	}
// 	return argsMap, nil
// }

// func createDependencyGraph(targets map[string]target.Target, trg target.Target, targetArgs []string) (dependency.Dependency, error) {
// 	var depConf target.DepConfig
// 	var err error
// 	depConf.ArgValues, err = parseCLITargetArgs(targetArgs)
// 	if err != nil {
// 		return dependency.Dependency{}, err
// 	}

// 	head, err := createDependencyFromTargets(trg, depConf, targets)
// 	if err != nil {
// 		return dependency.Dependency{}, err
// 	}
// 	return head, nil
// }
