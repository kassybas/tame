package steprunner

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
	"github.com/kassybas/tame/types/vartype"
)

func mergeOpts(globalOpts, targetOpts, stepOpts opts.ExecutionOpts) opts.ExecutionOpts {
	return opts.ExecutionOpts{
		Silent:  globalOpts.Silent || targetOpts.Silent || stepOpts.Silent,
		CanFail: globalOpts.CanFail || targetOpts.CanFail || stepOpts.CanFail,
		Async:   stepOpts.Async,
	}
}
func getIterableValues(iterableIf interface{}, vt *vartable.VarTable) ([]interface{}, error) {

	var iterableVal []interface{}
	switch iterableIf := iterableIf.(type) {
	case string:
		{
			iterable, err := vt.GetVar(iterableIf)
			if err != nil {
				return nil, fmt.Errorf("defined iterable cannot be resolved\n\t%s", err.Error())
			}
			if iterable.Type() != vartype.TListType && iterable.Type() != vartype.TMapType {
				return nil, fmt.Errorf("variable %s is not list or map (type: %T)", iterable.Name(), iterable)
			}
			var isList bool
			iterableVal, isList = iterable.Value().([]interface{})
			if !isList {
				iterableMap := iterable.Value().(map[interface{}]interface{})
				iterableVal = []interface{}{}
				for k := range iterableMap {
					iterableVal = append(iterableVal, k)
				}
			}
		}
	case []interface{}:
		{
			iterableVal = iterableIf
		}
	case map[interface{}]interface{}:
		{
			iterableVal = []interface{}{}
			for k := range iterableIf {
				iterableVal = append(iterableVal, k)
			}
		}
	default:
		{
			return nil, fmt.Errorf("unknown iterable")
		}
	}
	return iterableVal, nil
}

func getIters(vt *vartable.VarTable, s step.Step) (string, []interface{}, error) {
	if s.GetIteratorName() == "" && s.GetIterable() == nil {
		// No iterator and iterable -> no for loop, run once
		return "", []interface{}{""}, nil
	}
	// Iterable
	iterableIf := s.GetIterable()
	if iterableIf == nil {
		// nothing to iterate over -> run zero times
		return "", []interface{}{}, nil
	}
	iterableVal, err := getIterableValues(iterableIf, vt)
	if err != nil {
		return "", nil, err
	}
	// Iterator
	// validate iterator name
	if !strings.HasPrefix(s.GetIteratorName(), keywords.PrefixReference) {
		return "", nil, fmt.Errorf("iterator variable wrong format: %s (should be: %s%s)", s.GetIteratorName(), keywords.PrefixReference, s.GetIteratorName())
	}
	return s.GetIteratorName(), iterableVal, nil
}
func runStep(s step.Step, ctx tcontext.Context, vt *vartable.VarTable, parentOpts opts.ExecutionOpts) step.StepStatus {
	// Opts
	s.SetOpts(mergeOpts(ctx.Settings.GlobalOpts, parentOpts, s.GetOpts()))
	// to inherit the parent setting, we inject it in place of the global opts
	// TODO fix opt merging logic
	ctx.Settings.GlobalOpts = s.GetOpts()

	status := s.RunStep(ctx, vt)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("[step: %s]:: %s", s.GetName(), status.Err.Error())}
	}
	status.ResultNames = s.ResultNames()
	status.AllowedLessResults = s.Kind() == steptype.Shell
	// Breaking if it was breaking (return step) or the called step exec failed with non-zero exit
	status.IsBreaking = status.IsBreaking || (!s.GetOpts().CanFail && status.Stdstatus != 0)
	return status
}

func orchestrateIteration(iterator string, itVal interface{}, s step.Step, ctx tcontext.Context, vt *vartable.VarTable, wg *sync.WaitGroup, statusChan chan step.StepStatus, parentOpts opts.ExecutionOpts) step.StepStatus {
	vt.Add(iterator, itVal)
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
func startIterations(steps []step.Step, statusChan, resultChan chan step.StepStatus, syncStepDone chan bool, ctx tcontext.Context, vt *vartable.VarTable, parentOpts opts.ExecutionOpts) {
	var wg sync.WaitGroup
	for _, s := range steps {
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
				go orchestrateIteration(iterator, itVal, s, ctx, newVt, &wg, statusChan, parentOpts)
			} else {
				orchestrateIteration(iterator, itVal, s, ctx, vt, &wg, statusChan, parentOpts)
				// wait for sync step to finish processing results
				<-syncStepDone
			}
		}
	}
	// wait for all steps to finish
	wg.Wait()
	close(statusChan)
}

func RunAllSteps(steps []step.Step, ctx tcontext.Context, vt *vartable.VarTable, parentOpts opts.ExecutionOpts) step.StepStatus {
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
