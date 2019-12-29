package step

import (
	"github.com/kassybas/tame/types/opts"
	"github.com/kassybas/tame/types/steptype"
)

type BaseStep struct {
	name         string
	kind         steptype.Steptype
	resultNames  []string
	opts         opts.ExecutionOpts
	iteratorName string
	iterableName string
}
