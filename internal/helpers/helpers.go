package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/types/opts"
)

func GetKeyValueFromEnvString(envStr string) (string, string, error) {
	if !strings.Contains(envStr, "=") {
		return "", "", fmt.Errorf(`unknown argument format provided: expected: "arg_name=arg_value", got: %s`, envStr)
	}
	sps := strings.SplitN(envStr, keywords.CliArgSeparator, 2)
	k := sps[0]
	v := sps[1]
	return k, v, nil
}

func BuildOpts(optsDef []string) (opts.ExecutionOpts, error) {
	opts := opts.ExecutionOpts{}
	for _, opt := range optsDef {
		switch opt {
		case keywords.OptSilent:
			{
				opts.Silent = true
			}
		case keywords.OptCanFail:
			{
				opts.CanFail = true
			}
		default:
			{
				return opts, fmt.Errorf("unknown option: %s", opt)
			}
		}
	}
	return opts, nil
}

// returns the index and variable name
func ParseIndex(name string) (int, string, error) {
	lBr := strings.Index(name, keywords.IndexingSeparatorL) + 1
	rBr := strings.Index(name, keywords.IndexingSeparatorR)
	index, err := strconv.Atoi(name[lBr:rBr])
	if err != nil {
		return 0, "", fmt.Errorf("not integer index: %s %s", name, name[lBr:rBr])
	}
	return index, name[0 : lBr-1], nil
}
