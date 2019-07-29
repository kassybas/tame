package tvar

type VariableI interface {
	Type() TVarType
	Name() string
	Value() interface{}
	ToInt() (int, error)
	ToStr() string
	ToEnvVars() []string
	IsScalar() bool
}

// type Variable struct {
// 	Name string
// 	// TODO: interface
// 	stringValue string
// 	intValue    int
// 	Type        TVarType
// }

func CreateVariable(name string, value interface{}) VariableI {
	var newVar VariableI
	switch value.(type) {
	case string:
		{
			newVar = TString{
				name:  name,
				value: value.(string),
			}
		}
	case TString:
		{
			ts := value.(TString)
			newVar = TString{
				name:  name,
				value: ts.value,
			}
		}
	case map[interface{}]interface{}:
		{
			newVar = CreateMap(name, value)
		}
	case TMap:
		{
			tm := value.(TMap)
			newVar = TMap{
				name:  name,
				value: tm.value,
			}
		}
	}
	return newVar
}

type TVarType int

const (
	TErrorType TVarType = iota
	TStringType
	TIntType
	TFloatType
	TListType
	TMapType
)
