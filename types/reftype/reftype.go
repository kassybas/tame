package reftype

type RefType int

const (
	Unset RefType = iota
	Literal
	VarName
	InnerRef
	Index
)
