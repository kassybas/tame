package tvar

type VariableI interface {
	Type() TVarType
	Name() string
	Value() interface{}
	ToInt() int
	ToStr() string
	ToEnvVars() []string
}

// type Variable struct {
// 	Name string
// 	// TODO: interface
// 	stringValue string
// 	intValue    int
// 	Type        TVarType
// }

func CreateVariable(name string, value interface{}) VariableI {
	switch value.(type) {
	case string:
		{
			return TString{
				name:  name,
				value: value.(string),
			}
		}
	case TString:
		{
			ts := value.(TString)
			return TString{
				name:  name,
				value: ts.value,
			}
		}
	}
	// }
	// case int:
	// 	{
	// 		return Variable{
	// 			Name: name,
	// 			// intValue: value.(int),
	// 			// Type:     INT,
	// 		}
	// 	}
	// }
	return nil
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
