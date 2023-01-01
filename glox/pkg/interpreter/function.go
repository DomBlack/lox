package interpreter

import (
	"fmt"

	"github.com/DomBlack/lox/glox/pkg/ast"
)

type FunctionType uint8

const (
	FTNone FunctionType = iota
	FTFunction
	FTInitializer
	FTMethod
)

type Function struct {
	declaration   *ast.FunctionStmt
	closure       *Environment
	isInitializer bool
}

var _ Callable = (*Function)(nil)

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) Call(interpreter *interpreter, arguments []any) (rtn any) {
	defer func() {
		if r := recover(); r != nil {
			if rtnValue, ok := r.(*Return); ok {
				if f.isInitializer {
					rtn = f.closure.GetAt(0, "this")
				} else {
					rtn = rtnValue.Value
				}
			} else {
				panic(r)
			}
		}
	}()

	env := f.closure.Scope()
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.Body, env)

	if f.isInitializer {
		return f.closure.GetAt(0, "this")
	}

	return nil
}

func (f *Function) Bind(instance *Instance) *Function {
	env := f.closure.Scope()
	env.Define("this", instance)
	return &Function{f.declaration, env, f.isInitializer}
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
