package compile

import (
	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/schema"
	"github.com/kassybas/tame/types/settings"
)

func BuildSettings(tfs schema.SettingsShema) (settings.Settings, error) {
	opts, err := helpers.BuildOpts(tfs.GlobalOpts)
	if err != nil {
		return settings.Settings{}, err
	}
	if tfs.ShellFieldSeparator == "" {
		tfs.ShellFieldSeparator = keywords.ShellFieldSeparator
	}
	s := settings.Settings{
		UsedShell:           tfs.Shell,
		InitScript:          tfs.Init,
		GlobalOpts:          opts,
		ShieldEnv:           tfs.ShieldEnv,
		ShellFieldSeparator: tfs.ShellFieldSeparator,
	}

	return s, nil
}
