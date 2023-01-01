package token

import (
	"fmt"
)

type Token struct {
	Type    Type
	Lexeme  string
	Literal any
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s: %s %v", t.Type, t.Lexeme, t.Literal)
}
