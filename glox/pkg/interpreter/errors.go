package interpreter

import (
	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

func checkNumberOperands(operator *token.Token, operand ...any) {
	for _, o := range operand {
		if _, ok := o.(float64); !ok {
			panic(&errs.RuntimeError{Token: operator, Msg: "Operand must be a number."})
		}
	}
}
