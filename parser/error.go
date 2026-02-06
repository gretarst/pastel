package parser

import "fmt"

type ParserError struct {
	Msg    string
	Detail string
	Hint   string
	Line   int
	Column int
}

func (e *ParserError) Error() string {
	var msg string
	if e.Line > 0 {
		msg = fmt.Sprintf("\n[Parser Error] at line %d, column %d: %s", e.Line, e.Column, e.Msg)
	} else {
		msg = fmt.Sprintf("\n[Parser Error] %s", e.Msg)
	}
	if e.Detail != "" {
		msg += fmt.Sprintf("\n  → %s", e.Detail)
	}
	if e.Hint != "" {
		msg += fmt.Sprintf("\n  Hint: %s", e.Hint)
	}
	return msg
}
