package target

import (
	"fmt"
	"sync"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/steptype"
)

func (t Target) runStep(s step.Step, ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	// Check if condition
	resIf, err := evalConditionExpression(vt, s)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("[target: %s]:: %s", t.Name, err.Error())}
	}
	if !resIf {
		return step.StepStatus{}
	}
	// Opts
	s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, t.Opts, s.GetOpts()))
	// to inherit the parent setting, we inject it in place of the global opts
	// TODO fix opt merging logic
	ctx.Settings.GlobalOpts = s.GetOpts()

	status := s.RunStep(ctx, vt)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("[target: %s]:: %s", t.Name, status.Err.Error())}
	}
	status.ResultNames = s.ResultNames()
	status.AllowedLessResults = s.Kind() == steptype.Shell
	// Breaking if it was breaking (return step) or the called step exec failed with non-zero exit
	status.IsBreaking = status.IsBreaking || (!s.GetOpts().CanFail && status.Stdstatus != 0)
	return status
}

func (t *Target) orchestrateIteration(iterator string, itVal interface{}, s step.Step, ctx tcontext.Context, vt *vartable.VarTable, wg *sync.WaitGroup, statusChan chan step.StepStatus) step.StepStatus {
	vt.Add(iterator, itVal)
	status := t.runStep(s, ctx, vt)
	if status.Err != nil {
		status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
	}
	status.IsSync = !s.GetOpts().Async
	statusChan <- status
	wg.Done()
	return status
}

func processStatuses(statusChan, resultChan chan step.StepStatus, syncStepDone chan bool, vt *vartable.VarTable) {
	var lastStatus step.StepStatus
	for status := range statusChan {
		lastStatus = status
		if lastStatus.Err != nil {
			resultChan <- lastStatus
		}
		lastStatus.Err = updateVarsWithResultVariables(vt, status.ResultNames, status.Results, status.AllowedLessResults)
		if lastStatus.IsBreaking {
			resultChan <- lastStatus
		}
		if status.IsSync {
			syncStepDone <- true
		}
	}
	resultChan <- lastStatus
}

func (t *Target) startIterations(statusChan, resultChan chan step.StepStatus, syncStepDone chan bool, ctx tcontext.Context, vt *vartable.VarTable) {
	var wg sync.WaitGroup
	for _, s := range t.Steps {
		if s.Kind() == steptype.Wait {
			wg.Wait()
		}
		iterator, iterable, err := getIters(vt, s)
		if err != nil {
			resultChan <- step.StepStatus{Err: err, IsBreaking: true}
		}
		// if no for loop is defined then we iterate through one empty element
		for _, itVal := range iterable {
			wg.Add(1)
			if s.GetOpts().Async {
				newVt := vartable.CopyVarTable(vt)
				go t.orchestrateIteration(iterator, itVal, s, ctx, newVt, &wg, statusChan)
			} else {
				t.orchestrateIteration(iterator, itVal, s, ctx, vt, &wg, statusChan)
				// wait for sync step to finish processing results
				<-syncStepDone
			}
		}
	}
	// wait for all steps to finish
	wg.Wait()
	close(statusChan)
}

func (t *Target) runAllSteps(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	statusChan := make(chan step.StepStatus)
	resultChan := make(chan step.StepStatus)
	syncStepDone := make(chan bool, 1)
	// Start reading results
	go processStatuses(statusChan, resultChan, syncStepDone, vt)
	// Start creating iterations
	go t.startIterations(statusChan, resultChan, syncStepDone, ctx, vt)

	status := <-resultChan
	// setting it to false so it does not break the parent execution
	status.IsBreaking = false
	return status
}
