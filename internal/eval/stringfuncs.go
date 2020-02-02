package eval

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kassybas/tame/internal/helpers"
)

func Split(in, sep string) []string {
	// the last line would become and empty string in split so we trim it
	s := strings.TrimSuffix(in, "\n")
	return strings.Split(s, sep)
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
func CheckAtoi(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func ParseBool(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}
func CheckParseBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

func ParseFloat(s string, bitSize int) float64 {
	v, _ := strconv.ParseFloat(s, bitSize)
	return v
}
func CheckParseFloat(s string, bitSize int) bool {
	_, err := strconv.ParseFloat(s, bitSize)
	return err == nil
}

func Format(fmtString string, vals interface{}) string {
	vInterSlice, err := helpers.ConvertSliceToInterfaceSlice(vals)
	if err != nil {
		return fmt.Sprintf(fmtString, vals)
	}
	return fmt.Sprintf(fmtString, vInterSlice...)
}
