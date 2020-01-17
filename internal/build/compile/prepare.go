package compile

import (
	"fmt"
	"path/filepath"

	"github.com/kassybas/tame/schema"

	"github.com/kassybas/tame/internal/build/loader"
	"github.com/kassybas/tame/internal/build/targetparse"
	"github.com/kassybas/tame/internal/build/varparse"
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/step"
	"github.com/kassybas/tame/internal/tcontext"
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/settings"
	"github.com/sirupsen/logrus"
)

func createContext(globals []tvar.TVariable, sts settings.Settings) (tcontext.Context, error) {
	return tcontext.Context{
		Globals:  globals,
		Settings: sts,
	}, nil
}

func convertIncludesToRelativePath(path string, includes []schema.IncludeSchema) []schema.IncludeSchema {
	for i := range includes {
		includes[i].Path = fmt.Sprintf("%s%s%s", filepath.Dir(path), string(filepath.Separator), includes[i].Path)
	}
	return includes
}

func PrepareStep(path, targetName string, targetArgs map[string]interface{}) (step.Step, tcontext.Context, error) {
	// load static keys
	tf, dynamicKeys, err := loader.Load(path)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("error loading tamefile: %s\n\t%s", path, err.Error())
	}
	// load dynamic keys: targets
	if tf.Targets, err = targetparse.ParseTargets(dynamicKeys); err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("error while parsing targets in file %s\n\t%s", path, err.Error())
	}
	// load dynamic keys: global variables
	if tf.Globals, err = varparse.ParseGlobals(dynamicKeys); err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("error while parsing global variables in file: %s\n\t%s", path, err.Error())
	}
	if len(dynamicKeys) != 0 {
		return nil, tcontext.Context{}, fmt.Errorf("unknown keys in file: %s\n\t%s", path, err.Error())
	}
	tf.Includes = convertIncludesToRelativePath(path, tf.Includes)
	if !helpers.IsPublic(targetName) {
		return nil, tcontext.Context{}, fmt.Errorf("calling non-public target: %s\npublic targets must start with uppercase letter", targetName)
	}
	globals, err := varparse.EvaluateGlobals(tf.Globals)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("failed to evaluate global variables:\n\t%s", err.Error())
	}
	stgs, err := BuildSettings(tf.Sets)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("failed to evaluate settings:\n\t%s", err.Error())
	}
	ctx, err := createContext(globals, stgs)
	if err != nil {
		return nil, tcontext.Context{}, fmt.Errorf("error while creating context:\n\t%s", err.Error())
	}
	root, err := CompileTarget(tf, targetName, targetArgs, &ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	return root, ctx, nil
}
