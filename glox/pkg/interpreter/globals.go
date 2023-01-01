package interpreter

import (
	"time"
)

var globals = NewEnvironment()

func init() {
	globals.Define("clock", &CallableFunc{
		arity: 0,
		fn: func(interpreter Interpreter, arguments []any) any {
			return float64(time.Now().UnixNano()) / 1e9
		},
	})
}
