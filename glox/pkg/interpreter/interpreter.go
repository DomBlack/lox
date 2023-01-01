package interpreter

import (
	"fmt"

	"github.com/DomBlack/lox/glox/pkg/ast"
	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

type interpreter struct {
	environment *Environment
	locals      map[ast.Expr]int
}

type Interpreter interface {
	Resolve(stmts []ast.Stmt)
	Interpret(stmts []ast.Stmt)
}

func New() Interpreter {
	return &interpreter{
		environment: globals,
		locals:      make(map[ast.Expr]int),
	}
}

func (i *interpreter) Interpret(stmts []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*errs.RuntimeError); ok {
				errs.ErrorAtRuntime(e)
			} else {
				panic(fmt.Sprintf("Unhandled Panic (%T): %v", e, e))
			}
		}
	}()

	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *interpreter) Resolve(stmts []ast.Stmt) {
	resolver := &resolver{i, nil, FTNone, CTNone}
	resolver.resolve(stmts)
}

func (i *interpreter) resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *interpreter) lookupVariable(name *token.Token, expr ast.Expr) any {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	}

	return globals.Get(name)
}
