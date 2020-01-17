package compile

import (
	"os"

	"github.com/kassybas/tame/schema"

	"github.com/kassybas/tame/internal/build/targetbuild"
	"github.com/kassybas/tame/internal/helpscreen"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tcontext"
)

func getRootStepSchema(targetName string, cliVarArgs map[string]interface{}) (schema.MergedStepSchema, error) {
	var root schema.MergedStepSchema
	root.CalledTargetName = targetName
	root.CallArgumentsPassed = cliVarArgs
	return root, nil
}

func createDependencyGraph(targets map[string]target.Target, targetName string, cliVarArgs map[string]interface{}, includes []schema.IncludeSchema) (step.Step, error) {
	rootSchema, err := getRootStepSchema(targetName, cliVarArgs)
	if err != nil {
		return &callstep.CallStep{}, err
	}
	rootStep, err := callstep.NewCallStep(rootSchema)
	if err != nil {
		return &callstep.CallStep{}, err
	}
	calledTarget, err := findCalledTarget(targetName, targets, includes)
	if err != nil {
		return &callstep.CallStep{}, err
	}
	err = linkCalledTargets(&calledTarget.Steps, "[tame]", targets, includes)
	rootStep.SetCalledTarget(calledTarget)

	return rootStep, err
}

// CompileTarget creates the internal representation
func CompileTarget(tf schema.Tamefile, targetName string, cliVarArgs map[string]interface{}, ctx *tcontext.Context) (step.Step, error) {

	parsedTargets, err := targetbuild.BuildTargets(tf, ctx)
	if err != nil {
		return nil, err
	}
	if targetName == "" {
		helpscreen.PrintTeafileDescription(parsedTargets, tf)
		os.Exit(0)
	}
	// build the dependency graph with the called target
	var root step.Step
	root, err = createDependencyGraph(parsedTargets, targetName, cliVarArgs, tf.Includes)
	if err != nil {
		return root, err
	}

	return root, nil
}
