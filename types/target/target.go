package target

import (
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/settings"
	"github.com/kassybas/mate/types/step"
)

type ParamConfig struct {
	Name         string
	HasDefault   bool
	DefaultValue string
}
type Target struct {
	GlobalSettings *settings.Settings

	Name      string
	Body      string
	Return    []string
	Steps     []step.Step
	Params    []ParamConfig
	Opts      opts.ExecutionOpts
	Variables []Variable
}
type Variable struct {
	Name string
	// TODO: interface
	Value string
}
