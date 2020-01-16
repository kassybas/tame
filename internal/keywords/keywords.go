package keywords

var (
	PrefixReference = "$"

	OptSilent  = "silent"
	OptCanFail = "allowed-fail"
	OptAsync   = "async"
	OptsNotSet = "not-set"

	PossibleOpts = []string{
		OptCanFail, OptSilent, OptAsync,
	}

	ShellFieldSeparator = "_"
	TameFieldSeparator  = "."
	IndexingSeparatorL  = "["
	IndexingSeparatorR  = "]"

	CliArgSeparator        = "="
	GlobalDefaultVarSuffix = "?"

	StepIf = "if"
)
