package settings

import "github.com/kassybas/tame/types/opts"

type Settings struct {
	UsedShell           string
	InitScript          string
	GlobalOpts          opts.ExecutionOpts
	ShieldEnv           bool
	ShellFieldSeparator string
}
