package callstep

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func NewStep(map[string]interface{}) {
	var ok mapstructure.DecodeHookFunc
	fmt.Println(ok)
}
