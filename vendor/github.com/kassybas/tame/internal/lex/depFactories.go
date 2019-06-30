package lex

import (
	"fmt"
	"github.com/kassybas/tame/internal/dependency"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/target"
)

func findCalledTarget(name string, targets map[string]target.Target) (target.Target, error) {
	v, exists := targets[name]
	if exists {
		return v, nil
	}
	return target.Target{}, fmt.Errorf("Target not found: '%s'", name)
}


func getExecutionOpts(depConf target.DepConfig, callerExecOpts dependency.ExecutionOpts, settings target.Settings) (dependency.ExecutionOpts, error) {
	var execOpts dependency.ExecutionOpts
	var opts [] string
	opts = settings.DefaultOpts
	opts = append(opts, depConf.Opts...)

	for _, s := range opts {
		switch s {
		case "":
			// explicit empty (set.opts: "")
			continue
		case keywords.OptSilent:
			execOpts.Silent = true

		case keywords.OptStdRc:
			execOpts.SaveRc = true

		case keywords.OptStdout:
			execOpts.SaveOut = true

		case keywords.OptStderr:
			execOpts.SaveErr = true

		case keywords.OptOnce:
			execOpts.OnceIsEnough = true

		case keywords.OptParallel:
			execOpts.Parallel = true

		case keywords.OptCanFail:
			execOpts.CanFail = true

		default:
			return execOpts, fmt.Errorf("invalid execution option given: dependency '%s', opt '%s'", depConf.Name, s)
		}
	}

	// propagated opts from caller
	// silent is the only opt currently propagated
	if callerExecOpts.Silent {
		execOpts.Silent = true
	}
	return execOpts, nil
}

func createDependency(trg target.Target, depConf target.DepConfig, caller *dependency.Dependency, targets map[string]target.Target) (dependency.Dependency, error) {
	newDep := dependency.Dependency{
		Name:      trg.Name,
		ArgValues: depConf.ArgValues,
		Target:    trg,
		Caller:    caller,
	}

	var err error
	var callerExecOpts dependency.ExecutionOpts
	if caller == nil {
		callerExecOpts = dependency.ExecutionOpts{}
		newDep.GlobalSuccessfulDeps = &dependency.DependencyPool{}
	} else{
		callerExecOpts = caller.ExecOpts
		newDep.GlobalSuccessfulDeps = caller.GlobalSuccessfulDeps
	}
	newDep.ExecOpts, err = getExecutionOpts(depConf, callerExecOpts, *trg.GlobalSettings)


	if err != nil {
		return newDep, err
	}

	// TODO: overrides
	for _, depDepConfig := range newDep.Target.Deps {
		var newDepDep dependency.Dependency
		depDepTargetConfig, err := findCalledTarget(depDepConfig.Name, targets)
		if err != nil {
			return dependency.Dependency{}, err
		}
		newDepDep, err = createDependency(depDepTargetConfig, depDepConfig, &newDep, targets)
		if err != nil {
			return dependency.Dependency{}, err
		}
		newDep.Deps = append(newDep.Deps, newDepDep)
	}
	// TODO: post deps
	return newDep, err
}

func createDependencyFromTargets(trg target.Target, depConf target.DepConfig, targets map[string]target.Target) (dependency.Dependency, error) {

	dep, err := createDependency(trg, depConf, nil, targets)
	if err != nil {
		return dependency.Dependency{}, fmt.Errorf("could not create dependency graph:\nerror: %s \ncaller target: %s \ndependency: %v\n", err, trg.Name, depConf.Name)
	}
	return dep, err
}
