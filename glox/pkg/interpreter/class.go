package interpreter

import (
	"fmt"

	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

type ClassType uint8

const (
	CTNone ClassType = iota
	CTClass
	CTSubclass
)

type Class struct {
	Name    string
	Super   *Class
	Methods map[string]*Function
}

var _ Callable = (*Class)(nil)

func (c *Class) Arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.Arity()
}

func (c *Class) Call(interpreter *interpreter, arguments []any) any {
	instance := &Instance{Class: c}

	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (c *Class) findMethod(name string) *Function {
	if method, ok := c.Methods[name]; ok {
		return method
	}

	if c.Super != nil {
		return c.Super.findMethod(name)
	}

	return nil
}

func (c *Class) String() string {
	return c.Name
}

type Instance struct {
	Class  *Class
	Fields map[string]any
}

func (i *Instance) Get(name *token.Token) any {
	if value, ok := i.Fields[name.Lexeme]; ok {
		return value
	}

	if method := i.Class.findMethod(name.Lexeme); method != nil {
		return method.Bind(i)
	}

	panic(&errs.RuntimeError{Token: name, Msg: fmt.Sprintf("Undefined property '%s'.", name.Lexeme)})
}

func (i *Instance) Set(name *token.Token, value any) {
	i.Fields[name.Lexeme] = value
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s instance", i.Class.Name)
}
