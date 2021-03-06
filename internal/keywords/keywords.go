package keywords

var (
	PrefixReference = "$"

	OptSilent  = "silent"
	OptCanFail = "allow-fail"
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
	AppendToList        = "+"
	IndexingSeparatorL  = "["
	IndexingSeparatorR  = "]"

	CliArgSeparator        = "="
	GlobalDefaultVarSuffix = "?"

	IfStepPrefix = "if "
	IfStepElse   = "else"
	ShellStep    = "sh"
)
