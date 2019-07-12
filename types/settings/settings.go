package settings

import "github.com/kassybas/mate/types/opts"

type Settings struct {
	UsedShell  string
	InitScript string
	GlobalOpts opts.ExecutionOpts
	ShieldEnv  bool
}
