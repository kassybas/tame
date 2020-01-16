package helpers

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/opts"
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

func ConvertInterToMapStrInter(inter interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	mInter, ok := inter.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("converting non-map to map[string]: %v", inter)
	}
	for key, value := range mInter {
		switch key := key.(type) {
		case string:
			res[key] = value
		default:
			return nil, fmt.Errorf("non-string key found in map: %v", key)
		}
	}
	return res, nil
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
		case keywords.OptAsync:
			{
				opts.Async = true
			}
		default:
			{
				return opts, fmt.Errorf("unknown option: %s", opt)
			}
		}
	}
	return opts, nil
}
