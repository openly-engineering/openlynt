package lint

import (
	"go/token"
	"strings"
)

type Violation struct {
	Pos      token.Pos
	Position token.Position
	File     string
	Rule     *Rule

	Message string
}

func (e *Violation) Error() string {
	return e.Message
}

func (e *Violation) Line() int        { return e.Position.Line }
func (e *Violation) FilePath() string { return e.File }

type Violations struct {
	Violations []*Violation
}

func (e *Violations) Error() string {
	s := make([]string, len(e.Violations))
	for i := range e.Violations {
		s[i] = e.Violations[i].Error()
	}

	return strings.Join(s, "; ")
}
