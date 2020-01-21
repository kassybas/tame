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
		// case '"':
		// 	if singleQuotesStarted {
		// 		continue
		// 	}
		// 	if pos == prevPos+1 {
		// 		return nil, fmt.Errorf("quoting in the middle of variable reference: %s", fullName)
		// 	}
		// 	doubleQuotesStarted = !doubleQuotesStarted
		// case '\'':
		// 	if doubleQuotesStarted {
		// 		continue
		// 	}
		// 	if pos == prevPos+1 {
		// 		return nil, fmt.Errorf("quoting in the middle of variable reference: %s", fullName)
		// 	}
		// 	singleQuotesStarted = !singleQuotesStarted
		case '.':
			{
				if singleQuotesStarted || doubleQuotesStarted {
					continue
				}
				if pos != prevPos+1 {
					result.AddField(fullName[prevPos : pos-1])
				}
				prevPos = pos
			}
		case '[':
			{
				if singleQuotesStarted || doubleQuotesStarted {
					continue
				}
				if pos != prevPos+1 {
					result.AddField(fullName[prevPos : pos-1])
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
					result.AddField(fullName[prevPos : pos-1])
				}
				result.CloseInner()
				startedBrackets--
				prevPos = pos
				// TODO: check started and closed brackets
			}
		}
	}
	if prevPos != pos {
		result.AddField(fullName[prevPos:])
	}
	return result.CreateResultFields(), nil
}
