package keywords

var (
	Opts = "opts"

	PrefixReference   = "$"
	PrefixTameKeyword = "."

	OptsSeparator = " "
	OptSilent     = "silent"
	OptOnce       = "run-once"
	OptCanFail    = "allowed-fail"
	OptParallel   = "parallel"
	OptStdout     = "out"
	OptStderr     = "err"
	OptStdRc      = "rc"

	OptsNotSet = "not-set"

	PossibleOpts = []string{
		OptCanFail, OptSilent, OptOnce, OptParallel,
	}

	ShellFieldSeparator = "_"
	TameFieldSeparator  = "."
	IndexingSeparatorL  = "["
	IndexingSeparatorR  = "]"

	CliArgSeparator        = "="
	GlobalDefaultVarSuffix = "?"

	StepShell  = "sh"
	StepVar    = "var"
	StepCall   = "call"
	StepReturn = "return"

	StepFor         = "for"
	StepForIterable = "in"
	StepForIterator = "$"

	StepCallResult = "$"
	ShellOutResult = "$"
)
