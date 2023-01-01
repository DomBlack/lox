package interpreter

type Callable interface {
	Arity() int
	Call(interpreter *interpreter, arguments []any) any
}

type CallableFunc struct {
	arity int
	fn    func(interpreter Interpreter, arguments []any) any
}

var _ Callable = (*CallableFunc)(nil)

func (c *CallableFunc) Arity() int {
	return c.arity
}

func (c *CallableFunc) Call(interpreter *interpreter, arguments []any) any {
	return c.fn(interpreter, arguments)
}
