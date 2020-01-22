package dotref

import "fmt"

func ParseVarRef(fullName string) ([]RefField, error) {
	// var cur string
	if len(fullName) == 0 {
		return nil, fmt.Errorf("empty variable reference")
	}
	prevPos := 0
	pos := 0
	singleQuotesStarted := false
	doubleQuotesStarted := false
	var startedBrackets int
	tree := NewRefTree(nil)
	for _, ch := range fullName {
		pos++
		switch ch {
		case '"':
			if singleQuotesStarted {
				continue
			}
			doubleQuotesStarted = !doubleQuotesStarted
		case '\'':
			if doubleQuotesStarted {
				continue
			}
			singleQuotesStarted = !singleQuotesStarted
		case '.':
			{
				if singleQuotesStarted || doubleQuotesStarted {
					continue
				}
				if pos != prevPos+1 {
					err := tree.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return nil, err
					}
				}
				prevPos = pos
			}
		case '[':
			{
				if singleQuotesStarted || doubleQuotesStarted {
					continue
				}
				if pos != prevPos+1 {
					err := tree.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return nil, err
					}
				}
				tree.OpenInner()
				startedBrackets++
				prevPos = pos
			}
		case ']':
			{
				if singleQuotesStarted || doubleQuotesStarted {
					continue
				}
				if pos != prevPos+1 {
					err := tree.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return nil, err
					}
				}
				tree.CloseInner()
				startedBrackets--
				prevPos = pos
			}
		}
	}
	if startedBrackets != 0 {
		return nil, fmt.Errorf("unclosed brackets in variable reference: %s", fullName)
	}
	if prevPos != pos {
		tree.AddField(fullName[prevPos:])
	}
	allFields := tree.CreateResultFields()
	if allFields[0].FieldName == "" {
		return nil, fmt.Errorf("illegal variable reference: %s", fullName)
	}
	return allFields, nil
}
