package helpers

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/types/opts"
)

func GetKeyValueFromEnvString(envStr string) (string, interface{}, error) {
	// TODO: fix this with some proper regex
	if !strings.HasPrefix(envStr, "--") {
		return "", "", fmt.Errorf(`unknown argument format provided: expected: "--arg_name=arg_value", got: %s`, envStr)
	}
	// remove dashes
	envStr = strings.TrimPrefix(envStr, "--")
	// prepend $
	envStr = fmt.Sprintf("%s%s", keywords.PrefixReference, envStr)
	if !strings.Contains(envStr, keywords.CliArgSeparator) {
		return envStr, true, nil
	}

	sps := strings.SplitN(envStr, keywords.CliArgSeparator, 2)
	return sps[0], sps[1], nil
}

func ParseCLITargetArgs(flags []string) (map[string]interface{}, error) {
	args := make(map[string]interface{}, len(flags))
	for i, argStr := range flags {
		if i == 0 {
			continue
		}
		k, v, err := GetKeyValueFromEnvString(argStr)
		if err != nil {
			return nil, err
		}
		args[k] = v
	}
	return args, nil
}

func IsPublic(targetName string) bool {
	if len(targetName) == 0 {
		// helpscreen call -> TODO: fix this
		return true
	}
	firstChar := string(targetName[0])
	// check if first character is lowercase
	if strings.ToLower(firstChar) == firstChar {
		return false
	}
	return true
}

func DeepConvertInterToMapStrInter(inter interface{}) (interface{}, error) {
	var err error
	res := make(map[string]interface{})
	mInter, isMap := inter.(map[interface{}]interface{})
	if isMap {
		for key, value := range mInter {
			_, needToGoDeeper := value.(map[interface{}]interface{})
			if needToGoDeeper {
				value, err = DeepConvertInterToMapStrInter(value)
				if err != nil {
					return nil, err
				}
			}
			res[fmt.Sprintf("%v", key)] = value
			// return nil, fmt.Errorf("non-string key in map: %v", key)
		}
		return res, nil
	}
	// convert list elements, because list elements can be map[interface{}]interface{}
	lInter, isList := inter.([]interface{})
	if isList {
		for i := range lInter {
			lInter[i], err = DeepConvertInterToMapStrInter(lInter[i])
			if err != nil {
				return nil, err
			}
		}
		return lInter, nil
	}
	return inter, nil
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

func TrimLiteralQuotes(field string) (string, error) {
	if strings.HasPrefix(field, `"`) {
		if !strings.HasSuffix(field, `"`) {
			return "", fmt.Errorf("missing closing bracket: %s", field)
		}
		field = strings.Trim(field, `"`)
	} else if strings.HasPrefix(field, `'`) {
		if !strings.HasSuffix(field, `"`) {
			return "", fmt.Errorf("missing closing bracket: %s", field)
		}
		field = strings.Trim(field, `'`)
	}
	return field, nil
}
