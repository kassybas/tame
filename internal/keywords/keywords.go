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

	OptsDefaultValues = []string{OptStdRc, OptStderr, OptStdout}

	PossibleOpts = []string{
		OptCanFail, OptSilent, OptStderr, OptStdout, OptStdRc, OptOnce, OptParallel,
	}

	PrefixOut = ""
	PrefixErr = "err_"
	PrefixRc  = "rc_"

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

	StepCallResult    = "$"
	ShellOutResult    = "$"
	ShellErrResult    = "err$"
	ShellStatusResult = "status$"
)
