package steprunner

import (
	"fmt"
	"sync"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/stepblock"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
		Async:   stepOpts.Async,
	}
}

func runStep(s step.Step, ctx tcontext.Context, vt *vartable.VarTable, parentOpts opts.ExecutionOpts) step.StepStatus {
	// Opts
	s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, parentOpts, s.GetOpts()))
	// to inherit the parent setting, we inject it in place of the global opts
	// TODO fix opt merging logic
	ctx.Settings.GlobalOpts = s.GetOpts()

	status := s.RunStep(ctx, vt)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("[step: %s]:\n\t%s", s.GetName(), status.Err.Error())}
	}
	status.ResultNames = s.ResultNames()
	status.AllowedLessResults = s.Kind() == steptype.Shell
	// Breaking if it was breaking (return step) or the called step exec failed with non-zero exit
	status.IsBreaking = status.IsBreaking || (!s.GetOpts().CanFail && status.Stdstatus != 0)
	return status
}

func orchestrateIteration(s step.Step, ctx tcontext.Context, vt *vartable.VarTable, wg *sync.WaitGroup, statusChan chan step.StepStatus, parentOpts opts.ExecutionOpts) step.StepStatus {
	status := runStep(s, ctx, vt, parentOpts)
	if status.Err != nil {
		status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
	}
	status.IsSync = !s.GetOpts().Async
	statusChan <- status
	wg.Done()
	return status
}

func updateVarsWithResultVariables(vt *vartable.VarTable, resultVarNames []string, resultValues []interface{}, allowedLessResults bool) error {
	if len(resultVarNames) == 0 {
		return nil
	}
	if len(resultVarNames) > len(resultValues) {
		return fmt.Errorf("too many results expected, too little returned: %d > %d", len(resultVarNames), len(resultValues))
	}
	if len(resultVarNames) != len(resultValues) && !allowedLessResults {
		return fmt.Errorf("return and result variables do not match: %d != %d", len(resultVarNames), len(resultValues))
	}
	err := vt.Append(resultVarNames, resultValues)
	return err
}

func processStatuses(statusChan, resultChan chan step.StepStatus, syncStepDone chan bool, vt *vartable.VarTable) {
	var curStatus step.StepStatus
	for status := range statusChan {
		curStatus = status
		if curStatus.Err != nil {
			resultChan <- curStatus
		}
		curStatus.Err = updateVarsWithResultVariables(vt, status.ResultNames, status.Results, status.AllowedLessResults)
		if curStatus.IsBreaking {
			resultChan <- curStatus
		}
		if status.IsSync {
			syncStepDone <- true
		}
	}
	resultChan <- curStatus
}
func startIterations(steps stepblock.StepBlock, statusChan, resultChan chan step.StepStatus, syncStepDone chan bool, ctx tcontext.Context, vt *vartable.VarTable, parentOpts opts.ExecutionOpts) {
	var wg sync.WaitGroup
	for _, s := range steps.GetAll() {
		if s.Kind() == steptype.Wait {
			wg.Wait()
		}
		wg.Add(1)
		if s.GetOpts().Async {
			// copy vartable in case of async execution to have unique iterator
			newVt := vartable.CopyVarTable(vt)
			if s.GetIteratorVar() != nil {
				newVt.AddVar(s.GetIteratorVar())
			}
			go orchestrateIteration(s, ctx, newVt, &wg, statusChan, parentOpts)
		} else {
			vt.AddVar(s.GetIteratorVar())
			orchestrateIteration(s, ctx, vt, &wg, statusChan, parentOpts)
			// wait for sync step to finish processing results
			<-syncStepDone
		}
	}
	// wait for all steps to finish
	wg.Wait()
	close(statusChan)
}

func RunAllSteps(steps stepblock.StepBlock, ctx tcontext.Context, vt *vartable.VarTable, parentOpts opts.ExecutionOpts) step.StepStatus {
	statusChan := make(chan step.StepStatus)
	resultChan := make(chan step.StepStatus)
	syncStepDone := make(chan bool, 1)
	// Start reading results
	go processStatuses(statusChan, resultChan, syncStepDone, vt)
	// Start creating iterations
	go startIterations(steps, statusChan, resultChan, syncStepDone, ctx, vt, parentOpts)

	status := <-resultChan
	// setting it to false so it does not break the parent execution
	status.IsBreaking = false
	return status
}
