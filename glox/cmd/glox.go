package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/interpreter"
	"github.com/DomBlack/lox/glox/pkg/parser"
	"github.com/DomBlack/lox/glox/pkg/scanner"
)

var intpr = interpreter.New()

func main() {
	switch len(os.Args) {
	default:
		println("Usage: glox [script]")
		os.Exit(1)
	case 2:
		runFile(os.Args[1])
	case 1:
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read file: %v", err)
		os.Exit(2)
	}

	run(string(bytes))

	switch {
	case errs.HadError:
		os.Exit(65)
	case errs.HadRuntimeError:
		os.Exit(70)
	default:
		os.Exit(0)
	}
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			return
		}

		run(reader.Text())
		errs.HadError = false
	}
}

func run(source string) {
	s := scanner.New(source)
	tokens := s.ScanTokens()
	p := parser.New(tokens)
	stmts := p.Parse()

	if errs.HadError {
		return
	}

	intpr.Resolve(stmts)
	if errs.HadError {
		return
	}

	intpr.Interpret(stmts)
}
