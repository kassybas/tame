package dotref

func ParseDotRef(fullName string) ([]RefField, error) {
	// var cur string
	prevPos := 0
	pos := 0
	singleQuotesStarted := false
	doubleQuotesStarted := false
	var startedBrackets int
	result := NewRefTree(nil)
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
					err := result.AddField(fullName[prevPos : pos-1])
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
					err := result.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return nil, err
					}
				}
				result.OpenInner()
				startedBrackets++
				prevPos = pos
			}
		case ']':
			{
				if singleQuotesStarted || doubleQuotesStarted {
					continue
				}
				if pos != prevPos+1 {
					err := result.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return nil, err
					}
				}
				result.CloseInner()
				startedBrackets--
				prevPos = pos
			}
		}
	}
	if prevPos != pos {
		result.AddField(fullName[prevPos:])
	}
	return result.CreateResultFields(), nil
}
