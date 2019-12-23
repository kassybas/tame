package vartype

type TVarType int

const (
	TUnknownType TVarType = iota
	TScalarType
	TListType
	TMapType
)

func GetTypeNameString(t TVarType) string {
	switch t {
	case TScalarType:
		{
			return "TScalarType"
		}
	case TListType:
		{
			return "TListType"
		}
	case TMapType:
		{
			return "TMapType"
		}
	default:
		return "TUnknownType"
	}
}
