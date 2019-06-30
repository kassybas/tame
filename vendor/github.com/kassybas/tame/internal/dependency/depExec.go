package dependency

import (
	"github.com/kassybas/tame/internal/executor"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (td *Dependency) Exec() error {

	if td.ExecOpts.OnceIsEnough {
		prev, exists := td.FindPreviousSuccessfulExec()
		if exists{
			td.CopyResults(prev)
			return nil
		}
	}

	td.EvaluateReferences()
	executeAll(td.Deps)

	envVarStrings, err := td.GetAllEnvVars()
	if err != nil {
		return err
	}
	script := td.Target.GlobalSettings.InitScript + "\n" + td.Target.Script
	//scriptOut, scriptErr, scriptRc, err := executor.ExecuteScript(td.Name, script, envVarStrings, td.Target.GetShell(), td.ExecOpts.Silent, td.Target.GlobalSettings.ShieldEnv)
	scriptOut, scriptErr, scriptRc, err := executor.ExecuteScript(td.Name, script, envVarStrings, td.Target.GetShell(), td.ExecOpts.Silent, td.Target.GlobalSettings.ShieldEnv)
	if err != nil {
		logrus.Fatal("Finishing the script execution failed on a tame level:", err)
	}

	failed := scriptRc != "0"
	td.Executed = true
	td.Failed = failed

	if failed && !td.ExecOpts.CanFail {
		logrus.Error("Executing dependency: ", td.Name, " failed with exit status: ", scriptRc)
		if td.ExecOpts.Silent {
			logrus.Error(td.Name, " stdout:\n", scriptOut)
			logrus.Error(td.Name, " stderr:\n", scriptErr)
		}

		rc, err := strconv.Atoi(scriptRc)
		if err != nil {
			return err
		}
		// TODO: move this
		os.Exit(rc)
	}

	if td.ExecOpts.SaveOut {
		// Trim trailing newline
		scriptOut = strings.TrimSuffix(scriptOut, "\n")
		td.stdOut = scriptOut
	}
	if td.ExecOpts.SaveErr {
		// Trim trailing newline
		scriptErr = strings.TrimSuffix(scriptErr, "\n")
		td.stdErr = scriptErr
	}
	if td.ExecOpts.SaveRc {
		td.stdRc = scriptRc
	}
	err = executeAll(td.PostExecDeps)
	if err != nil {
		return err
	}
	// Optimize this to determine to only add when it will be needed again
	td.GlobalSuccessfulDeps.Deps = append(td.GlobalSuccessfulDeps.Deps, *td)

	return nil
}

func executeAll(deps []Dependency) error {
	var parallelDepPool []Dependency

	for i := range deps {
		if deps[i].Executed {
			// Skip once executed deps
			continue
		}
		if deps[i].ExecOpts.Parallel {
			// Collect parallel dependencies to start them at once
			parallelDepPool = append(parallelDepPool, deps[i])
			continue
		}
		if parallelDepPool != nil {
			// In case dep is not parallel, execute the previously collected parallel pool (if not empty)
			err := executeDepsParallel(parallelDepPool)
			if err != nil {
				return err
			}
			// Empty parallel pool
			parallelDepPool = nil
		}
		err := deps[i].Exec()
		if err != nil {
			return err
		}
	}
	if parallelDepPool != nil {
		// If parallel pool is not empty, execute the rest
		return executeDepsParallel(parallelDepPool)
	}
	return nil
}

func executeDepsParallel(deps []Dependency) error {
	var wg sync.WaitGroup
	var err error
	for i := range deps {
		wg.Add(1)
		xi := i
		go func() {
			locErr := deps[xi].Exec()
			wg.Done()
			if locErr != nil {
				err = locErr
			}
		}()
	}
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}
