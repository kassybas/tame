package lex

import (
	"github.com/kassybas/mate/internal/helpers"
	"github.com/kassybas/mate/schema"
	"github.com/kassybas/mate/types/settings"
)

func BuildSettings(tfs schema.SettingsDefintion) (settings.Settings, error) {
	opts, err := helpers.BuildOpts(tfs.GlobalOpts)
	if err != nil {
		return settings.Settings{}, err
	}

	s := settings.Settings{
		UsedShell:  tfs.Shell,
		InitScript: tfs.Init,
		GlobalOpts: opts,
		ShieldEnv:  tfs.ShieldEnv,
	}

	return s, nil
}
