package ast

import (
	"testing"

	"github.com/DomBlack/lox/glox/pkg/token"
)

func TestPrinter(t *testing.T) {
	expression := &Binary{
		Left: &Unary{
			Operator: &token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
			Right: &Literal{
				Value: 123,
			},
		},
		Operator: &token.Token{Type: token.STAR, Lexeme: "*", Line: 1},
		Right: &Grouping{
			Expression: &Literal{
				Value: 45.67,
			},
		},
	}

	printer := &printer{}
	println(printer.Print(expression))
}
