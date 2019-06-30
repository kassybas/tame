package target

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type DepConfig struct {
	Name      string
	Opts      []string
	ArgValues map[string]string
}

type ParamConfig struct {
	Name         string
	HasDefault   bool
	DefaultValue string
	Value        string
}

type Target struct {
	GlobalSettings *Settings
	GlobalVars     *[]GlobalVar
	Name           string
	Script         string
	Deps           []DepConfig
	Params         []ParamConfig
}

func (t *Target) GetDefaultValue(paramName string) (string, error) {
	for _, p := range t.Params {
		if p.Name == paramName {
			if p.HasDefault {
				return p.DefaultValue, nil
			}
			return "", fmt.Errorf("default not set or value passed for: %s", p.Name)
		}
	}
	return "", fmt.Errorf("parameter does not exist for target: '%s', arg: '%s'", t.Name, paramName)
}

func (t *Target) GetShell() string {
	var sh string
	sh = os.Getenv("SHELL")

	if t.GlobalSettings.UsedShell != "" {
		sh = t.GlobalSettings.UsedShell
	}

	if sh == "" {
		logrus.Warn("Shell could not be determined: falling back to /bin/sh. Please set SHELL to environment or use the 'set shell /path/to/shell' directive")
		sh = "/bin/sh"
	}
	return sh
}
