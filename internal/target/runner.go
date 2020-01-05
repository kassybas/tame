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

func (t *Target) orchestrateStep(s step.Step, ctx tcontext.Context, vt *vartable.VarTable, wg *sync.WaitGroup, doneChan chan bool, statusChan chan step.StepStatus, isLast bool) step.StepStatus {
	var status step.StepStatus
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
				status = t.runStep(sLoc, ctx, &newVt)
				if status.Err != nil {
					status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
				}
				statusChan <- status
				if isLast {
					close(statusChan)
				}
			}()
		} else {
			vt.Add(iterator, itVal)
			status = t.runStep(s, ctx, vt)
			if status.Err != nil {
				status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
			}
			statusChan <- status
			if isLast {
				close(statusChan)
			}
		}
	}
	if !s.GetOpts().Async {
		wg.Done() // sync execution wait for finishing at the end
	}
	return status
}

func (t *Target) runAllSteps(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	statusChan := make(chan step.StepStatus)
	doneChan := make(chan bool, 1)
	var lastStatus step.StepStatus
	// Start reading results
	go func() {
		for status := range statusChan {
			lastStatus = status
			if status.Err != nil || status.IsBreaking {
				// if it failed we are closing the done channel
				doneChan <- true
			}
			err := updateVarsWithResultVariables(vt, status.ResultNames, status.Results, status.AllowedLessResults)
			if err != nil {
				lastStatus = step.StepStatus{Err: err}
				doneChan <- true
			}
		}
		doneChan <- true
	}()
	// Run steps
	var wg sync.WaitGroup
	for i, s := range t.Steps {
		wg.Wait() // waits if the execution should be sync
		if !s.GetOpts().Async {
			wg.Add(1)
		}
		go t.orchestrateStep(s, ctx, vt, &wg, doneChan, statusChan, i == len(t.Steps)-1)
	}
	<-doneChan
	// setting it to false so it does not break the parent execution
	lastStatus.IsBreaking = false
	return lastStatus
}
