package exprtype

type ExprType int

const (
	Unset ExprType = iota
	Literal
	VarName
	InnerRef
	Index
)
