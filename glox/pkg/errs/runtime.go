package errs

import (
	"fmt"
	"os"

	"github.com/DomBlack/lox/glox/pkg/token"
)

type RuntimeError struct {
	Token *token.Token
	Msg   string
}

func ErrorAtRuntime(e *RuntimeError) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", e.Msg, e.Token.Line)
	HadRuntimeError = true
}
