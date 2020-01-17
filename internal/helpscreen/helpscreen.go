package helpscreen

import (
	"fmt"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/target"
)

func PrintTeafileDescription(targets map[string]target.Target) {
	fmt.Println("Available targets:")
	for k, v := range targets {
		if !helpers.IsPublic(k) {
			continue
		}
		fmt.Printf("\t- %s", k)
		if len(v.Summary) != 0 {
			fmt.Printf("-- %s", v.Summary)
		}
		if len(v.Params) != 0 {
			fmt.Printf("\n\t\t\\\\-args: ")
			for _, p := range v.Params {
				fmt.Printf("%s", p.Name)
				if p.HasDefault {
					fmt.Printf(" (default: '%s') ", p.DefaultValue)
				}
			}
		}
		fmt.Printf("\n")
	}
}
