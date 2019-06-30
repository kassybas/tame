package keywords

var (
	Arg    = "args"
	Opts   = "opts"
	Script = "script"

	PrefixReference = "$"

	OptsSeparator = " "
	OptSilent     = "silent"
	OptOnce       = "run-once"
	OptCanFail    = "allowed-fail"
	OptParallel   = "parallel"
	OptStdout     = "out"
	OptStderr     = "err"
	OptStdRc      = "rc"

	OptsNotSet = "not-set"


	OptsDefaultValues = []string{OptStdRc, OptStderr, OptStdout}

	PossibleOpts = []string{
		OptCanFail, OptSilent, OptStderr, OptStdout, OptStdRc, OptOnce, OptParallel,
	}

	PrefixOut = ""
	PrefixErr = "err_"
	PrefixRc  = "rc_"

	CliArgSeparator = "="
	GlobalDefaultVarSuffix = "?"
)
