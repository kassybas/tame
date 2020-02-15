package eval

import (
	"strconv"
	"strings"
)

func AllFunctions() map[string]interface{} {
	return map[string]interface{}{
		// strings
		"stringsSplit":        Split, // special fix required, see at function definition
		"stringsContainsAny":  strings.ContainsAny,
		"stringsCount":        strings.Count,
		"stringsEqualFold":    strings.EqualFold,
		"stringsFields":       strings.Fields,
		"stringsHasPrefix":    strings.HasPrefix,
		"stringsHasSuffix":    strings.HasSuffix,
		"stringsIndex":        strings.Index,
		"stringsIndexAny":     strings.IndexAny,
		"stringsLastIndex":    strings.LastIndex,
		"stringsLastIndexAny": strings.LastIndexAny,
		"stringsRepeat":       strings.Repeat,
		"stringsReplace":      strings.Replace,
		"stringsReplaceAll":   strings.ReplaceAll,
		"stringsSplitAfter":   strings.SplitAfter,
		"stringsSplitAfterN":  strings.SplitAfterN,
		"stringsSplitN":       strings.SplitN,
		"stringsTitle":        strings.Title,
		"stringsToLower":      strings.ToLower,
		"stringsToTitle":      strings.ToTitle,
		"stringsToUpper":      strings.ToUpper,
		"stringsToValidUTF8":  strings.ToValidUTF8,
		"stringsTrim":         strings.Trim,
		"stringsTrimLeft":     strings.TrimLeft,
		"stringsTrimPrefix":   strings.TrimPrefix,
		"stringsTrimRight":    strings.TrimRight,
		"stringsTrimSpace":    stringsTrimSpace,
		"stringsTrimSuffix":   strings.TrimSuffix,
		// strconv
		"strconvAtoi":            Atoi,
		"strconvCheckAtoi":       CheckAtoi,
		"strconvFormatBool":      strconv.FormatBool,
		"strconvItoa":            strconv.Itoa,
		"strconvParseBool":       ParseBool,
		"strconvCheckParseBool":  CheckParseBool,
		"strconvParseFloat":      ParseFloat,
		"strconvCheckParseFloat": CheckParseFloat,
		// builtins
		"append": Append,
		"extend": Extend,
		"format": Format,
	}
}
