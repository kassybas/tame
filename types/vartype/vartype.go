package vartype

type TVarType int

const (
	TUnknownType TVarType = iota
	TScalarType
	TListType
	TMapType
)

func (t TVarType) Name() string {
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
