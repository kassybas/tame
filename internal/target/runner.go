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

func (t *Target) orchestrateIteration(iterator string, itVal interface{}, s step.Step, ctx tcontext.Context, vt *vartable.VarTable, wg *sync.WaitGroup, statusChan chan step.StepStatus) step.StepStatus {
	vt.Add(iterator, itVal)
	status := t.runStep(s, ctx, vt)
	if status.Err != nil {
		status.Err = fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())
	}
	statusChan <- status
	wg.Done()
	return status
}

func processStatuses(statusChan chan step.StepStatus, resultChan chan step.StepStatus, vt *vartable.VarTable) {
	var lastStatus step.StepStatus
	for status := range statusChan {
		lastStatus = status
		lastStatus.Err = updateVarsWithResultVariables(vt, status.ResultNames, status.Results, status.AllowedLessResults)
		fmt.Println("--WE RUN")
		if lastStatus.IsBreaking {
			fmt.Println("--WE BROKE")
			resultChan <- lastStatus
		}
	}
	fmt.Println("--WE SET THE RESULT")
	resultChan <- lastStatus
	fmt.Println("--WE ARE DONE SETTING THE RESULT")
}

func (t *Target) runAllSteps(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	statusChan := make(chan step.StepStatus)
	resultChan := make(chan step.StepStatus, 1)
	var wg sync.WaitGroup
	// Start reading results
	go processStatuses(statusChan, resultChan, vt)
	// Run steps
	for _, s := range t.Steps {
		iterator, iterable, err := getIters(vt, s)
		if err != nil {
			return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)}
		}
		// if no for loop is defined
		// iterable is one empty element
		// iterator is empty string which is ignored during adding to the vartable
		for _, itVal := range iterable {
			select {
			case status := <-resultChan:
				{
					// setting it to false so it does not break the parent execution
					fmt.Println("--WE RECEIVED THE RESULT")
					status.IsBreaking = false
					return status
				}
			default:
				{
					wg.Add(1)
					var newVt *vartable.VarTable
					if s.GetOpts().Async {
						newVt = vartable.CopyVarTable(vt)
					} else {
						newVt = vt
					}
					go t.orchestrateIteration(iterator, itVal, s, ctx, newVt, &wg, statusChan)
					// if !s.GetOpts().Async {
					// 	wg.Wait() // waits if the execution should be sync
					// }
				}
			}
		}
	}
	wg.Wait()
	close(statusChan)
	// should never get here
	return step.StepStatus{Err: fmt.Errorf("internal error in parallel execution")}
}
