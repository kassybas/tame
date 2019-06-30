package parser

import (
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/internal/tamefile"
	"os"
	"strings"
)

func buildTargetDefinition(targetKey string, targetContainer tamefile.TargetContainer, settings *target.Settings, globals *[]target.GlobalVar) (target.Target, error) {
	newTarget := target.Target{
		Name:           targetKey,
		Script:         targetContainer.Script,
		GlobalSettings: settings,
		GlobalVars:     globals,
	}

	// Arguments
	for _, argValue := range targetContainer.ArgContainer {
		arg, err := buildArgumentDefinition(argValue)
		if err != nil {
			return newTarget, err
		}
		newTarget.Params = append(newTarget.Params, arg)
	}

	for _, depValue := range targetContainer.DepContainer {
		dep, err := buildDependencyDefinition(depValue)
		if err != nil {
			return newTarget, err
		}
		newTarget.Deps = append(newTarget.Deps, dep)
	}
	return newTarget, nil
}

func buildSettingsDefinition(tfs tamefile.SetConfig) (target.Settings, error) {
	var settings target.Settings
	settings.UsedShell = tfs.ShellContainer
	settings.InitScript = tfs.InitContainer
	settings.ShieldEnv = tfs.ShieldEnvContainer
	if tfs.DefaultOptsContainer == keywords.OptsNotSet {
		settings.DefaultOpts = keywords.OptsDefaultValues
	} else {
		settings.DefaultOpts = strings.Split(tfs.DefaultOptsContainer, keywords.OptsSeparator)
	}

	return settings, nil
}

func getGlobalVar(name, value string)(target.GlobalVar, error){
	var g target.GlobalVar
	if strings.HasSuffix(name, keywords.GlobalDefaultVarSuffix){
		trimmedVarName := strings.TrimSuffix(name, keywords.GlobalDefaultVarSuffix)
		g.Name = strings.TrimSpace(trimmedVarName)

		sysEnvValue, sysEnvExists := os.LookupEnv(helpers.FlattenEnvVarNameGlobal(g.Name))
		if sysEnvExists {
			g.Value = sysEnvValue
		}else {
			g.Value = value
		}
	} else{
		g.Name = name
		g.Value = value
	}
	g.EnvVarName = helpers.FlattenEnvVarNameGlobal(g.Name)
	return g, nil
}

func buildGlobalVariables(tfg map[string]string) ([]target.GlobalVar, error) {
	var globals []target.GlobalVar
	for k, v := range tfg {
		g, err := getGlobalVar(k,v)
		if err!=nil{
			return nil, err
		}
		globals = append(globals, g)
	}
	return globals, nil
}

func ParseTeafile(tf tamefile.Teafile) (map[string]target.Target, error) {

	settings, err := buildSettingsDefinition(tf.Sets)
	if err != nil {
		return nil, err
	}
	globals, err := buildGlobalVariables(tf.Globals)

	targets := make(map[string]target.Target)
	for targetKey, targetValue := range tf.Targets {
		trg, err := buildTargetDefinition(targetKey, targetValue, &settings, &globals)
		if err != nil {
			return targets, err
		}
		targets[targetKey] = trg
	}

	return targets, nil
}
