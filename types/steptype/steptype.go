package steptype

type Steptype int

const (
	Unset Steptype = iota
	Call
	Shell
	Var
	Return
	Expr
	Wait
	If
	For
	Dump
	Load
	Print
)

func (t Steptype) ToStr() string {
	switch t {
	case Unset:
		return "unset"
	case Call:
		return "call"
	case Shell:
		return "shell"
	case Var:
		return "var"
	case Return:
		return "return"
	case Expr:
		return "expr"
	case Wait:
		return "wait"
	case If:
		return "if"
	case For:
		return "for"
	case Dump:
		return "dump"
	case Load:
		return "load"
	case Print:
		return "print"
	default:
		panic("internal error: non-steptype steptype")
	}

}
