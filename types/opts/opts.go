package opts

type ExecutionOpts struct {
	Silent       bool
	SaveOut      bool
	SaveErr      bool
	SaveRc       bool
	CanFail      bool
	OnceIsEnough bool
	Parallel     bool
}
