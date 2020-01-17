package compile

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/stepblock"

	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/step/forstep"
	"github.com/kassybas/tame/internal/step/ifstep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

func loadCalledTargetInclude(name string, includes []schema.IncludeSchema, targets map[string]target.Target) (target.Target, error) {
	namespace := strings.Split(name, keywords.TameFieldSeparator)[0]
	calledTargetName := strings.TrimLeft(name, namespace+keywords.TameFieldSeparator)
	for _, incl := range includes {
		if namespace == incl.Alias {
			s, _, err := PrepareStep(incl.Path, calledTargetName, map[string]interface{}{})
			if err != nil {
				return target.Target{}, fmt.Errorf("error while loading include: %s\n\t%s", name, err.Error())
			}
			return s.(*callstep.CallStep).GetCalledTarget(), err
		}
	}
	return target.Target{}, fmt.Errorf("namespace of referenced target not found in includes: no alias matching '%s'", namespace)
}

func findCalledTarget(name string, targets map[string]target.Target, includes []schema.IncludeSchema) (target.Target, error) {
	if strings.Contains(name, keywords.TameFieldSeparator) {
		return loadCalledTargetInclude(name, includes, targets)
	}
	v, exists := targets[name]
	if exists {
		return v, nil
	}
	return target.Target{}, fmt.Errorf("target not found: '%s'", name)
}

func linkCalledTargets(steps *stepblock.StepBlock, caller string, targets map[string]target.Target, includes []schema.IncludeSchema) error {
	for i := range steps.GetAll() {
		s := *steps.Get(i)
		switch s.Kind() {
		case steptype.Call:
			{
				calledTarget, err := findCalledTarget(s.GetName(), targets, includes)
				if err != nil {
					return fmt.Errorf("caller: '%s'\n\t%s", caller, err.Error())
				}
				// Including a target already populates it
				if s.(*callstep.CallStep).IsCalledTargetSet() {
					// setting called targets of child steps
					err = linkCalledTargets(&calledTarget.Steps, s.GetName(), targets, includes)
					if err != nil {
						return err
					}
				}
				s.(*callstep.CallStep).SetCalledTarget(calledTarget)
			}
		case steptype.If:
			{
				err := linkCalledTargets(s.(*ifstep.IfStep).GetIfSteps(), s.GetName(), targets, includes)
				if err != nil {
					return err
				}
				err = linkCalledTargets(s.(*ifstep.IfStep).GetElseSteps(), s.GetName(), targets, includes)
				if err != nil {
					return err
				}
			}
		case steptype.For:
			{
				err := linkCalledTargets(s.(*forstep.ForStep).GetForSteps(), s.GetName(), targets, includes)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
