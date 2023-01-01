package errs

import (
	"fmt"
	"os"

	"github.com/DomBlack/lox/glox/pkg/token"
)

var (
	HadError        bool
	HadRuntimeError bool
)

func ErrorOnLine(line int, message string) {
	report(line, "", message)
}

func ErrorAtToken(t *token.Token, message string) {
	if t.Type == token.EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func report(line int, where string, message string) {
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	HadError = true
}
