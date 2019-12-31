package parse

import (
	"fmt"
	"strings"
)

func parseVariableName(k string) (string, error) {
	fields := strings.Fields(k)
	if len(fields) > 2 {
		return "", fmt.Errorf("'%s': variable name contains whitespaces", k)
	}
	if len(fields) == 1 {
		return "", fmt.Errorf("'%s': no variable target name found: (correct: var $varname: value)", k)
	}
	return fields[1], validateVariableName(fields[1])
}
