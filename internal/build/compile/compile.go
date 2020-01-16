package compile

import (
	"os"

	"github.com/kassybas/tame/schema"

	"github.com/kassybas/tame/internal/build/targetparse"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/helpscreen"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/step/callstep"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tcontext"
)

func parseCLITargetArgs(targetArgs []string) (map[string]interface{}, error) {
	args := make(map[string]interface{}, len(targetArgs))
	for _, argStr := range targetArgs {
		k, v, err := helpers.GetKeyValueFromEnvString(argStr)
		if err != nil {
			return nil, err
		}
		args[k] = v
	}
	return args, nil
}

func getRootStepSchema(targetName string, cliVarArgs []string) (schema.MergedStepSchema, error) {
	var root schema.MergedStepSchema
	var err error
	root.CalledTargetName = &targetName
	root.CallArgumentsPassed, err = parseCLITargetArgs(cliVarArgs)
	if err != nil {
		return root, err
	}
	return root, err
}

func createDependencyGraph(targets map[string]target.Target, targetName string, cliVarArgs []string, includes []schema.IncludeSchema) (step.Step, error) {
	rootSchema, err := getRootStepSchema(targetName, cliVarArgs)
	if err != nil {
		return &callstep.CallStep{}, err
	}
	rootStep, err := callstep.NewCallStep(rootSchema)
	if err != nil {
		return &callstep.CallStep{}, err
	}
	calledTarget, err := findCalledTarget(targetName, "[tame cli]", targets, includes)
	if err != nil {
		return &callstep.CallStep{}, err
	}
	err = populateSteps(&calledTarget, targets, includes)
	rootStep.SetCalledTarget(calledTarget)

	return rootStep, err
}

// Compile creates the internal representation
func Compile(tf schema.Tamefile, targetName string, cliVarArgs []string, ctx *tcontext.Context) (step.Step, error) {

	parsedTargets, err := targetparse.ParseTeafile(tf, ctx)
	if err != nil {
		return nil, err
	}
	if targetName == "" {
		helpscreen.PrintTeafileDescription(parsedTargets)
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
