package dependency

import (
	"github.com/kassybas/tame/internal/target"
)

type ExecutionOpts struct {
	Silent       bool
	SaveOut      bool
	SaveErr      bool
	SaveRc       bool
	CanFail      bool
	OnceIsEnough bool
	Parallel     bool
}

type Dependency struct {
	Name     string
	Target   target.Target
	ExecOpts ExecutionOpts

	GlobalSuccessfulDeps *DependencyPool
	Caller               *Dependency
	ArgValues            map[string]string
	Deps                 []Dependency
	PostExecDeps         []Dependency

	stdOut string
	stdErr string
	stdRc  string

	Executed       bool
	Failed         bool
}

type DependencyPool struct {
	Deps []Dependency
}

func (td *Dependency) FindPreviousSuccessfulExec() (Dependency, bool){
	for _, d := range td.GlobalSuccessfulDeps.Deps {
		if td.Name == d.Name {
			return d, true
		}
	}
	return Dependency{}, false
}

func (td *Dependency) CopyResults(source Dependency){
	td.Executed = source.Executed
	td.stdRc = source.stdRc
	td.stdErr = source.stdErr
	td.Failed = source.Failed
	td.ArgValues = source.ArgValues
	td.ExecOpts = source.ExecOpts
}

