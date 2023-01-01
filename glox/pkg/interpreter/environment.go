package interpreter

import (
	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{Values: make(map[string]any)}
}

func (e *Environment) Scope() *Environment {
	return &Environment{Enclosing: e, Values: make(map[string]any)}
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Enclosing
	}

	return env
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name *token.Token) any {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	panic(&errs.RuntimeError{Token: name, Msg: "Undefined variable '" + name.Lexeme + "'."})
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.ancestor(distance).Values[name]
}

func (e *Environment) Assign(name *token.Token, value any) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	panic(&errs.RuntimeError{Token: name, Msg: "Undefined variable '" + name.Lexeme + "'."})
}

func (e *Environment) AssignAt(distance int, name *token.Token, value any) {
	e.ancestor(distance).Values[name.Lexeme] = value
}
