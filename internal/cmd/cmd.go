package cmd

import (
	"github.com/kassybas/mate/internal/orchestrate"
)

// MakeCommand runs the given target of the file
func MakeCommand(file, targetName string, targetArgs []string) {
	orchestrate.Make(file, targetName, targetArgs)
}
