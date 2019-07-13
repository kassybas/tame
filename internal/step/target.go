package step

import (
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/settings"
)

type Param struct {
	Name         string
	HasDefault   bool
	DefaultValue string
}
type Target struct {
	GlobalSettings *settings.Settings

	Name      string
	Return    []string
	Steps     []Step
	Params    []Param
	Opts      opts.ExecutionOpts
	Variables []tvar.Variable
	Summary   string
}
