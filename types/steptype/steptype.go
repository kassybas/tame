package steptype

type Steptype int

const (
	Unset Steptype = iota
	Call
	Shell
	Var
	Return
	Expr
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
	default:
		panic("internal error: non-steptype steptype")
	}

}
