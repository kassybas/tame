package helpscreen

import (
	"fmt"

	"github.com/kassybas/tame/internal/helpers"
	"github.com/kassybas/tame/internal/target"
	"github.com/kassybas/tame/schema"
)

func PrintTeafileDescription(targets map[string]target.Target, tf schema.Tamefile) {
	fmt.Println("\t[[Summary]]")
	if tf.Summary != "" {
		fmt.Printf("\t\t%s\n", tf.Summary)
	}
	fmt.Println("\t[[Targets]]")
	for k, v := range targets {
		if !helpers.IsPublic(k) {
			continue
		}
		fmt.Printf("\t\t- %s", k)
		if len(v.Summary) != 0 {
			fmt.Printf("-- %s", v.Summary)
		}
		if len(v.Params) != 0 {
			fmt.Printf("\n\t\t  |- args: ")
			for _, p := range v.Params {
				fmt.Printf("%s", p.Name)
				if p.HasDefault {
					fmt.Printf(" (default: '%v') ", p.DefaultValue)
				}
			}
		}
		fmt.Printf("\n")
	}
}
