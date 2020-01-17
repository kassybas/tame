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
	// illegal target names
	ReservedKeywords = []string{
		"settings", "return", "wait", "if", "else", "opts", "sh", "$", "include",
	}

	ShellFieldSeparator = "_"
	TameFieldSeparator  = "."
	IndexingSeparatorL  = "["
	IndexingSeparatorR  = "]"

	CliArgSeparator        = "="
	GlobalDefaultVarSuffix = "?"

	IfStepPrefix = "if "
	IfStepElse   = "else"
	ShellStep    = "sh"
)
