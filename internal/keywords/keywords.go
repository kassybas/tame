package keywords

var (
	Opts = "opts"

	PrefixReference   = "$"
	PrefixTameKeyword = "."

	OptsSeparator = " "
	OptSilent     = "silent"
	OptCanFail    = "allowed-fail"
	OptAsync      = "async"
	OptsNotSet    = "not-set"

	PossibleOpts = []string{
		OptCanFail, OptSilent, OptAsync,
	}

	ShellFieldSeparator = "_"
	TameFieldSeparator  = "."
	IndexingSeparatorL  = "["
	IndexingSeparatorR  = "]"

	CliArgSeparator        = "="
	GlobalDefaultVarSuffix = "?"

	StepShell  = "sh"
	StepReturn = "return"

	StepFor         = "for"
	StepForIterable = "in"
	StepForIterator = "$"

	StepCallResult = "$"
	ShellOutResult = "$"
)
