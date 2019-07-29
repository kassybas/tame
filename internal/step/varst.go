package step

import (
	"strings"

	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/tvar"
	"github.com/kassybas/mate/internal/vartable"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
)

type VarStep struct {
	Definitions map[string]interface{} // tvar.VariableI
	Opts        opts.ExecutionOpts
	Results     Result
}

func (s VarStep) GetName() string {
	keys := make([]string, 0, len(s.Definitions))
	for key := range s.Definitions {
		keys = append(keys, key)
	}
	return strings.Join(keys, ",")
}

func (s *VarStep) Kind() steptype.Steptype {
	return steptype.Var
}

func (s *VarStep) SetOpts(o opts.ExecutionOpts) {
	s.Opts = o
}

func (s *VarStep) GetResult() Result {
	return s.Results
}
func (s *VarStep) RunStep(ctx tcontext.Context, vt vartable.VarTable) error {
	// TODO: eval variables
	s.Results.ResultVars = make([]string, len(s.Definitions))
	s.Results.ResultValue = make([]tvar.VariableI, len(s.Definitions))
	i := 0
	for k, v := range s.Definitions {
		s.Results.ResultVars[i] = k
		s.Results.ResultValue[i] = tvar.CreateVariable(k, v)
		i++
	}
	return nil
}

func (s *VarStep) GetCalledTargetName() string {
	return s.GetName()
}

func (s *VarStep) GetOpts() opts.ExecutionOpts {
	return s.Opts
}

func (s *VarStep) SetCalledTarget(t Target) {

}
