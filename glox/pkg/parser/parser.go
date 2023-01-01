package parser

import (
	"github.com/DomBlack/lox/glox/pkg/ast"
	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() ast.Stmt {
	switch {
	case p.match(token.CLASS):
		return p.classDeclaration()
	case p.match(token.FUN):
		return p.function("function")
	case p.match(token.VAR):
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) classDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect class name.")

	var superclass *ast.VariableExpr
	if p.match(token.LESS) {
		p.consume(token.IDENTIFIER, "Expect superclass name.")
		superclass = &ast.VariableExpr{Name: p.previous()}
	}

	p.consume(token.LEFT_BRACE, "Expect '{' before class body.")

	var methods []*ast.FunctionStmt
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after class body.")
	return &ast.ClassStmt{Name: name, Superclass: superclass, Methods: methods}
}

func (p *Parser) function(kind string) *ast.FunctionStmt {
	name := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	var params []*token.Token
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			params = append(params, p.consume(token.IDENTIFIER, "Expect parameter name."))
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()

	return &ast.FunctionStmt{Name: name, Params: params, Body: body}
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return &ast.VarStmt{Name: name, Initializer: initializer}
}

func (p *Parser) statement() ast.Stmt {
	switch {
	case p.match(token.FOR):
		return p.forStatement()
	case p.match(token.IF):
		return p.ifStatement()
	case p.match(token.PRINT):
		return p.printStatement()
	case p.match(token.RETURN):
		return p.returnStatement()
	case p.match(token.WHILE):
		return p.whileStatement()
	case p.match(token.LEFT_BRACE):
		return &ast.BlockStmt{Statements: p.block()}
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &ast.BlockStmt{Statements: []ast.Stmt{
			body,
			&ast.ExpressionStmt{Expression: increment},
		}}
	}

	if condition == nil {
		condition = &ast.LiteralExpr{Value: true}
	}
	body = &ast.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.BlockStmt{Statements: []ast.Stmt{
			initializer,
			body,
		}}
	}

	return body
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() ast.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &ast.PrintStmt{Expression: value}
}

func (p *Parser) returnStatement() ast.Stmt {
	keyword := p.previous()
	var value ast.Expr
	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return &ast.ReturnStmt{Keyword: keyword, Value: value}
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return &ast.WhileStmt{Condition: condition, Body: body}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return &ast.ExpressionStmt{Expression: expr}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) expression() (rtn ast.Expr) {
	defer func() {
		if r := recover(); r != nil {
			if r == ErrParseError {
				rtn = nil
			} else {
				panic(r)
			}
		}
	}()

	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		switch v := expr.(type) {
		case *ast.VariableExpr:
			name := v.Name
			return &ast.AssignExpr{Name: name, Value: value}
		case *ast.GetExpr:
			return &ast.SetExpr{Object: v.Object, Name: v.Name, Value: value}
		default:
			errs.ErrorAtToken(equals, "Invalid assignment target.")
		}
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}

	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

paramsLoop:
	for {
		switch {
		case p.match(token.LEFT_PAREN):
			expr = p.finishCall(expr)
		case p.match(token.DOT):
			name := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			expr = &ast.GetExpr{Object: expr, Name: name}
		default:
			break paramsLoop
		}
	}

	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	var arguments []ast.Expr

	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				_ = p.error(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, p.expression())
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")

	return &ast.CallExpr{
		Callee:    callee,
		Arguments: arguments,
		Paren:     paren,
	}
}

func (p *Parser) primary() ast.Expr {
	switch {
	case p.match(token.FALSE):
		return &ast.LiteralExpr{Value: false}
	case p.match(token.TRUE):
		return &ast.LiteralExpr{Value: true}
	case p.match(token.NIL):
		return &ast.LiteralExpr{Value: nil}
	case p.match(token.NUMBER, token.STRING):
		return &ast.LiteralExpr{Value: p.previous().Literal}
	case p.match(token.LEFT_PAREN):
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.GroupingExpr{Expression: expr}
	case p.match(token.SUPER):
		keyword := p.previous()
		p.consume(token.DOT, "Expect '.' after 'super'.")
		method := p.consume(token.IDENTIFIER, "Expect superclass method name.")
		return &ast.SuperExpr{Keyword: keyword, Method: method}
	case p.match(token.THIS):
		return &ast.ThisExpr{Keyword: p.previous()}
	case p.match(token.IDENTIFIER):
		return &ast.VariableExpr{Name: p.previous()}
	default:
		panic(p.error(p.peek(), "Expect expression."))
	}
}

func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) consume(t token.Type, message string) *token.Token {
	if p.check(t) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(t *token.Token, message string) error {
	errs.ErrorAtToken(t, message)
	return ErrParseError
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS:
		case token.FUN:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}

		p.advance()
	}
}
