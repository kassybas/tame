package exprparse

import (
	"fmt"
)

type ParseState struct {
	inSingleQuotes  bool
	inDoubleQuotes  bool
	squareBrStarted int
	roundBrStarted  int
}

func (s *ParseState) InQuotes() bool {
	return s.inSingleQuotes || s.inDoubleQuotes
}

func (s *ParseState) InQuotesOrRoundBrackets() bool {
	return s.inSingleQuotes || s.inDoubleQuotes || s.roundBrStarted > 0
}

func (s *ParseState) InRoundBracket() bool {
	return s.roundBrStarted > 0
}

func (s *ParseState) CloseRoundBracket() error {
	s.roundBrStarted--
	if s.roundBrStarted < 0 {
		return fmt.Errorf("closing round bracket without opening")
	}
	return nil
}

func (s *ParseState) OpenRoundBracket() {
	s.roundBrStarted++
}

func ParseExpression(fullName string) (ParseTree, error) {
	// var cur string
	if len(fullName) == 0 {
		return ParseTree{}, fmt.Errorf("empty variable reference")
	}
	prevPos := 0
	pos := 0
	state := ParseState{}
	tree := NewRefTree(nil)
	for _, ch := range fullName {
		pos++
		switch ch {
		case '(':
			if state.InQuotes() {
				continue
			}
			state.OpenRoundBracket()
		case ')':
			if state.InQuotes() {
				continue
			}
			err := state.CloseRoundBracket()
			if err != nil {
				return ParseTree{}, fmt.Errorf("failed to parse expression: %s\n\t%s", fullName, err.Error())
			}
		case '"':
			if state.inSingleQuotes || state.InRoundBracket() {
				continue
			}
			state.inDoubleQuotes = !state.inDoubleQuotes
		case '\'':
			if state.inDoubleQuotes || state.InRoundBracket() {
				continue
			}
			state.inSingleQuotes = !state.inSingleQuotes
		case '.':
			{
				if state.InQuotesOrRoundBrackets() {
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
				if state.InQuotesOrRoundBrackets() {
					continue
				}
				if pos != prevPos+1 {
					err := tree.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return ParseTree{}, err
					}
				}
				tree.OpenInner()
				state.squareBrStarted++
				prevPos = pos
			}
		case ']':
			{
				if state.InQuotesOrRoundBrackets() {
					continue
				}
				if pos != prevPos+1 {
					err := tree.AddField(fullName[prevPos : pos-1])
					if err != nil {
						return ParseTree{}, err
					}
				}
				tree.CloseInner()
				state.squareBrStarted--
				prevPos = pos
			}
		}
	}
	if state.squareBrStarted != 0 || state.InRoundBracket() {
		return ParseTree{}, fmt.Errorf("unclosed brackets in expression: %s", fullName)
	}
	if prevPos != pos {
		tree.AddField(fullName[prevPos:])
	}
	tree.cur = nil
	tree.parent = nil
	return *tree, nil
}
