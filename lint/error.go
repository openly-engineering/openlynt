package lint

import (
	"go/token"
	"strings"
)

type Error struct {
	Pos      token.Pos
	Position token.Position

	Message string
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorCollection struct {
	Errors []*Error
}

func (e *ErrorCollection) Error() string {
	s := make([]string, len(e.Errors))
	for i := range e.Errors {
		s[i] = e.Errors[i].Error()
	}

	return strings.Join(s, "; ")
}
