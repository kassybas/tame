package tcontext

import (
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/settings"
)

type Context struct {
	Globals  []tvar.Variable
	Settings settings.Settings
}