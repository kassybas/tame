package vartype

type TVarType int

const (
	TUnknownType TVarType = iota
	TStringType
	TIntType
	TFloatType
	TListType
	TMapType
	TBoolType
	TNullType
)

func GetTypeNameString(t TVarType) string {
	switch t {
	case TStringType:
		{
			return "TStringType"
		}
	case TIntType:
		{
			return "TIntType"
		}
	case TFloatType:
		{
			return "TFloatType"
		}
	case TListType:
		{
			return "TListType"
		}
	case TMapType:
		{
			return "TMapType"
		}
	case TBoolType:
		{
			return "TBoolType"
		}
	case TNullType:
		{
			return "TNullType"
		}
	default:
		return "TUnknownType"
	}
}
