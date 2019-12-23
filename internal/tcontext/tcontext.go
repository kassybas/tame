package tcontext

import (
	"github.com/kassybas/tame/internal/tvar"
	"github.com/kassybas/tame/types/settings"
)

type Context struct {
	Globals  []tvar.TVariable
	Settings settings.Settings
}
