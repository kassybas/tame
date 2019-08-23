package step

import (
	"github.com/kassybas/mate/internal/tcontext"
	"github.com/kassybas/mate/internal/vartable"
	"github.com/kassybas/mate/types/opts"
	"github.com/kassybas/mate/types/steptype"
)

type VarStep struct {
	Definition interface{}
	Opts       opts.ExecutionOpts
	Results    Result
	Name       string
}

func (s VarStep) GetName() string {
	return s.Name
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
	s.Results.ResultNames = []string{s.Name}
	s.Results.ResultValues = []interface{}{s.Definition}
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
