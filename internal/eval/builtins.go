package eval

import "github.com/kassybas/tame/internal/helpers"

import "github.com/sirupsen/logrus"

func Append(data []interface{}, value interface{}) []interface{} {
	return append(data, value)
}

func Extend(data []interface{}, values interface{}) []interface{} {
	switch values.(type) {
	case []interface{}:
		return append(data, values.([]interface{})...)
	case []int, []string, []bool, []float64:
		list, err := helpers.ConvertSliceToInterfaceSlice(values)
		if err != nil {
			logrus.Fatalf("failed to extend '%v' to '%v'\n\t%s", data, values, err.Error())
		}
		return append(data, list...)
	default:
		return Append(data, values)
	}
}
