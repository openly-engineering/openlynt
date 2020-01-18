package lint

import "go/token"

type Error struct {
	Pos      token.Pos
	Position token.Position

	Message string
}

func (e *Error) Error() string {
	return e.Message
}
