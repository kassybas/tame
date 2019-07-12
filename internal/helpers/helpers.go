package helpers

import (
	"fmt"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/types/step"
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

func GetKeyValueFromEnvString(envStr string) (string, string, error) {
	if !strings.Contains(envStr, "=") {
		return "", "", fmt.Errorf(`unknown argument format provided: expected: "arg_name=arg_value", got: %s`, envStr)
	}
	sps := strings.SplitN(envStr, keywords.CliArgSeparator, 2)
	k := sps[0]
	v := sps[1]
	return k, v, nil
}

func FormatEnvVars(vars map[string]step.Variable) []string {
	formattedVars := []string{}
	for _, v := range vars {
		// Remove $ for shell env format
		trimmedName := strings.TrimPrefix(v.Name, keywords.PrefixReference)
		newVar := trimmedName + "=" + v.Value
		formattedVars = append(formattedVars, newVar)
	}
	return formattedVars
}
