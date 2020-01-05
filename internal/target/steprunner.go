package target

import (
	"fmt"

	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/types/steptype"
)

func (t Target) runStep(s step.Step, ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	// Check if condition
	resIf, err := evalConditionExpression(*vt, s)
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

	newVt := *vt
	status := s.RunStep(ctx, newVt)
	if status.Err != nil {
		return step.StepStatus{Err: fmt.Errorf("[target: %s]:: %s", t.Name, status.Err.Error())}
	}
	// Breaking if it was breaking (return step) or the called step exec failed with non-zero exit
	status.IsBreaking = status.IsBreaking || (!s.GetOpts().CanFail && status.Stdstatus != 0)
	err = updateVarsWithResultVariables(vt, s.ResultNames(), status.Results, s.Kind() == steptype.Shell)
	if err != nil {
		return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), err.Error())}
	}
	return status
}

func (t *Target) runAllSteps(ctx tcontext.Context, vt vartable.VarTable) step.StepStatus {
	var status step.StepStatus
	for _, s := range t.Steps {
		// TODO: refactor to more dry
		if s.GetIterable() == nil {
			status = t.runStep(s, ctx, &vt)
			if status.Err != nil {
				return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())}
			}
			if status.IsBreaking {
				// setting it to false so it does not break the parent execution
				status.IsBreaking = false
				return status
			}
		} else {
			iterator, iterable, err := getIters(vt, s)
			if err != nil {
				return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), err)}
			}
			for _, itVal := range iterable {
				vt.Add(iterator, itVal)
				status = t.runStep(s, ctx, &vt)
				if status.Err != nil {
					return step.StepStatus{Err: fmt.Errorf("in step: %s\n\t%s", s.GetName(), status.Err.Error())}
				}
				if status.IsBreaking {
					// setting it to false so it does not break the parent execution
					status.IsBreaking = false
					return status
				}
			}
		}
	}
	return status

}
