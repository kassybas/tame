package forstep

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/basestep"
	"github.com/kassybas/tame/internal/stepblock"
	"github.com/kassybas/tame/internal/steprunner"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/internal/vartable"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

type ForStep struct {
	basestep.BaseStep
	iteratorName string
	iterable     interface{}
	forSteps     stepblock.StepBlock
}

func NewForStep(stepDef schema.MergedStepSchema, forSteps []step.Step) (*ForStep, error) {
	var newStep ForStep
	var err error
	if stepDef.ForLoop == nil {
		return nil, fmt.Errorf("no iterators found in for loop definition: %s", stepDef.ForLoop)
	}
	if len(*stepDef.ForLoop) > 1 {
		return nil, fmt.Errorf("multiple iterators found in for loop definition: %s", stepDef.ForLoop)
	}
	for k, v := range *stepDef.ForLoop {
		newStep.iteratorName = k
		newStep.iterable = v
	}
	newStep.forSteps = stepblock.NewStepBlock(forSteps)
	newStep.BaseStep, err = basestep.NewBaseStep(stepDef, steptype.For, fmt.Sprintf("for %s: %v", newStep.iteratorName, newStep.iterable))
	return &newStep, err
}

func (s *ForStep) getIterableValues(vt *vartable.VarTable) ([]interface{}, error) {
	var iterableVal []interface{}
	resIterable, err := vt.ResolveValue(s.iterable)
	if err != nil {
		return nil, fmt.Errorf("could not resolve iterable expression: %s\n\t%s", s.iterable, err.Error())
	}
	switch iterableV := resIterable.(type) {
	case []interface{}:
		{
			iterableVal = iterableV
		}
	case []int, []string, []float64, []bool:
		{
			iterableVal, err = helpers.ConvertSliceToInterfaceSlice(resIterable)
			if err != nil {
				return nil, fmt.Errorf("could not determine iterable: %s\n\t%s", s.iterable, err.Error())
			}
		}
	case map[interface{}]interface{}:
		{
			iterableVal = []interface{}{}
			// in map we iterate through the keys
			for k := range iterableV {
				iterableVal = append(iterableVal, k)
			}
		}
	case map[string]interface{}:
		{
			iterableVal = []interface{}{}
			// in map we iterate through the keys
			for k := range iterableV {
				iterableVal = append(iterableVal, k)
			}
		}
	default:
		{
			return nil, fmt.Errorf("unknown iterable type: %v (type: %T)", resIterable, resIterable)
		}
	}
	return iterableVal, nil
}
func (s *ForStep) getIters(vt *vartable.VarTable) (string, []interface{}, error) {
	// Iterable
	if s.iterable == nil {
		// nothing to iterate over -> run zero times
		return "", []interface{}{}, nil
	}
	iterableVal, err := s.getIterableValues(vt)
	if err != nil {
		return "", nil, err
	}
	// Iterator
	// validate iterator name
	if !strings.HasPrefix(s.iteratorName, keywords.PrefixReference) {
		return "", nil, fmt.Errorf("iterator variable wrong format: %s (should be: %s%s)", s.iteratorName, keywords.PrefixReference, s.iteratorName)
	}
	return s.iteratorName, iterableVal, nil
}

func (s *ForStep) GetForSteps() *stepblock.StepBlock {
	return &s.forSteps
}

func (s *ForStep) RunStep(ctx tcontext.Context, vt *vartable.VarTable) step.StepStatus {
	iterator, iterable, err := s.getIters(vt)
	if err != nil {
		return step.StepStatus{
			Err: fmt.Errorf("could not determine iteration in step: '%s'\n\t%s", s.GetName(), err.Error()),
		}
	}
	var status step.StepStatus

	// generate flattened list of steps
	genForSteps := []step.Step{}
	for _, itVal := range iterable {
		for _, fStep := range s.forSteps.GetAll() {
			newStep := step.Clone(fStep)
			newStep.SetIteratorVar(tvar.NewVariable(iterator, itVal))
			genForSteps = append(genForSteps, newStep)
		}
	}
	status = steprunner.RunAllSteps(stepblock.NewStepBlock(genForSteps), ctx, vt, s.GetOpts())
	if status.Err != nil {
		return status
	}
	return status
}
