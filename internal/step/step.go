package step

import (
	"strings"

	"github.com/kassybas/mate/internal/vartable"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
)

type StepI interface {
	GetName() string
	Kind() steptype.Steptype
	GetResult() Result
	GetOpts() opts.ExecutionOpts
	SetOpts(opts.ExecutionOpts)
	RunStep(tcontext.Context, vartable.VarTable) error
	GetCalledTargetName() string
	SetCalledTarget(Target)
}

type Result struct {
	StdoutVar    string
	StdoutValue  string
	StderrVar    string
	StderrValue  string
	StdrcVar     string
	StdrcValue   int
	ResultVars   []string
	ResultValues []string
}

func FormatEnvVars(vars map[string]tvar.Variable) []string {
	formattedVars := []string{}
	for _, v := range vars {
		// Remove $ for shell env format
		trimmedName := strings.TrimPrefix(v.Name, keywords.PrefixReference)
		newVar := trimmedName + "=" + v.Value
		formattedVars = append(formattedVars, newVar)
	}
	return formattedVars
}
