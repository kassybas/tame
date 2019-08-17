package helpscreen

import (
	"fmt"

	"github.com/kassybas/mate/internal/step"
)

func PrintTeafileDescription(targets map[string]step.Target) {
	fmt.Println("Available targets:")
	for k, v := range targets {
		fmt.Printf("- %s", k)
		if v.Params == nil {
			fmt.Println()
			continue
		}
		fmt.Print(":\n")
		fmt.Printf("    @params: ")

		for _, p := range v.Params {
			fmt.Printf("%s", p.Name)
			if p.HasDefault {
				fmt.Printf(" (default: %s), ", p.DefaultValue)
			}
		}
		fmt.Printf("\n")
	}
}
