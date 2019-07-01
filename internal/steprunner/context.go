package steprunner

import (
	"github.com/kassybas/mate/types/settings"
	"github.com/kassybas/mate/types/step"
)

type Context struct {
	Globals  []step.Variable
	Settings settings.Settings
}
