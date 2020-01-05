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

func (t *Target) orchestrateStep(s step.Step, ctx tcontext.Context, vt *vartable.VarTable, wg *sync.WaitGroup, doneChan chan bool, statusChan chan step.StepStatus) step.StepStatus {
	iterator, iterable, err := getIters(vt, s)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)}
	}

	// if no for loop is defined
	// iterable is one empty element
	// iterator is empty string which is ignored during adding to the vartable
	for _, itVal := range iterable {
		if s.GetOpts().Async {
			sLoc := s
			newVt := vartable.CopyVarTable(vt)
			newVt.Add(iterator, itVal)
			go func() {
				status := t.runStep(sLoc, ctx, &newVt)
				if status.Err != nil {
					status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
				}
				statusChan <- status
				if status.IsBreaking {
					doneChan <- true
				}
			}()
		} else {
			vt.Add(iterator, itVal)
			status := t.runStep(s, ctx, vt)
			if status.Err != nil {
				status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
			}
			statusChan <- status
			if status.IsBreaking {
				doneChan <- true
			}
		}
	}
	if !s.GetOpts().Async {
		wg.Done() // sync execution wait for finishing at the end
	}
	return step.StepStatus{}
}

func (t *Target) runAllSteps(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	statusChan := make(chan step.StepStatus)
	resultChan := make(chan step.StepStatus, 1)
	doneChan := make(chan bool, 1)
	var wg sync.WaitGroup
	// Start reading results
	go func() {
		var lastStatus step.StepStatus
		for status := range statusChan {
			lastStatus = status
			lastStatus.Err = updateVarsWithResultVariables(vt, status.ResultNames, status.Results, status.AllowedLessResults)
			if lastStatus.IsBreaking {
				resultChan <- lastStatus
				return
			}
		}
		resultChan <- lastStatus
	}()
	taskChan := make(chan step.Step)
	go func() {
		for _, s := range t.Steps {
			taskChan <- s
		}
		doneChan <- true
	}()
	// Run steps
	for {
		select {
		case <-doneChan:
			{
				status := <-resultChan
				status.IsBreaking = false
				return status
			}
		case stepTask := <-taskChan:
			{
				if !stepTask.GetOpts().Async {
					wg.Add(1)
				}
				go t.orchestrateStep(stepTask, ctx, vt, &wg, doneChan, statusChan)
				if !stepTask.GetOpts().Async {
					wg.Wait() // waits if the execution should be sync
				}
			}
		}
	}
	// setting it to false so it does not break the parent execution
}
