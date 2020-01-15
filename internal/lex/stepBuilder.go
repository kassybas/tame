package lex

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"

	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/steptype"
)

func loadCalledTargetInclude(name, caller string, includes []schema.IncludeSchema, targets map[string]target.Target) (target.Target, error) {
	namespace := strings.Split(name, keywords.TameFieldSeparator)[0]
	calledTargetName := strings.TrimLeft(name, namespace+keywords.TameFieldSeparator)
	for _, incl := range includes {
		if namespace == incl.Alias {
			s, _, err := PrepareStep(incl.Path, calledTargetName, []string{})
			if err != nil {
				return target.Target{}, fmt.Errorf("error while loading include: %s\n\t%s", name, err.Error())
			}
			return s.(*callstep.CallStep).GetCalledTarget(), err
		}
	}
	return target.Target{}, fmt.Errorf("namespace of referenced target not found in includes: no alias matching '%s' [called by: '%s']", namespace, caller)
}

func findCalledTarget(name, caller string, targets map[string]target.Target, includes []schema.IncludeSchema) (target.Target, error) {
	if strings.Contains(name, keywords.TameFieldSeparator) {
		return loadCalledTargetInclude(name, caller, includes, targets)
	}
	v, exists := targets[name]
	if exists {
		return v, nil
	}
	return target.Target{}, fmt.Errorf("target not found: '%s' [called by: '%s']", name, caller)
}

func populateSteps(trg *target.Target, targets map[string]target.Target, includes []schema.IncludeSchema) error {
	for i := range trg.Steps {
		if trg.Steps[i].Kind() == steptype.Call {
			calledTarget, err := findCalledTarget(trg.Steps[i].GetName(), trg.Name, targets, includes)
			if err != nil {
				return err
			}
			// Including a target already populates it, we only populate steps if
			// they were not already populated
			if calledTarget.Name == "" {
				err = populateSteps(&calledTarget, targets, includes)
				if err != nil {
					return err
				}
			}
			trg.Steps[i].(*callstep.CallStep).SetCalledTarget(calledTarget)
		}
	}
	return nil
}
