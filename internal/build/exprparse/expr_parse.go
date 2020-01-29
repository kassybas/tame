package exprparse

import (
	"fmt"
)

func ParseExpression(fullName string) (ParseTree, error) {
	// var cur string
	if len(fullName) == 0 {
		return ParseTree{}, fmt.Errorf("empty variable reference")
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
						return ParseTree{}, err
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
						return ParseTree{}, err
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
						return ParseTree{}, err
					}
				}
				tree.CloseInner()
				startedBrackets--
				prevPos = pos
			}
		}
	}
	if startedBrackets != 0 {
		return ParseTree{}, fmt.Errorf("unclosed brackets in expression: %s", fullName)
	}
	if prevPos != pos {
		tree.AddField(fullName[prevPos:])
	}
	tree.cur = nil
	tree.parent = nil
	return *tree, nil
}
