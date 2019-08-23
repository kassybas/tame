package steptype

type Steptype int

const (
	Unset Steptype = iota
	Call
	Shell
	Var
	Return
)
