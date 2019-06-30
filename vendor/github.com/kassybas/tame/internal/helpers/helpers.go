package helpers

import (
	"fmt"
	"github.com/kassybas/tame/internal/keywords"
	"strings"
)

func FlattenEnvVarNameLocal(name string) string{
	name = strings.Replace(name, "-", "_", -1)
	name = strings.Replace(name, ".", "_", -1)
	return name
}

func FlattenEnvVarNameGlobal(name string) string{
	name = strings.Replace(name, "-", "_", -1)
	name = strings.Replace(name, ".", "_", -1)
	return name
}


func GetKeyValueFromEnvString(envStr string)(string, string, error){
	if ! strings.Contains(envStr, "="){
		return "","", fmt.Errorf(`unknown argument format provided: expected: "arg-name=arg-value", got: %s`, envStr)
	}
	sps := strings.SplitN(envStr, keywords.CliArgSeparator, 2)
	k := sps[0]
	v := sps[1]
	return k, v, nil
}