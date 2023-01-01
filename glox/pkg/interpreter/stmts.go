package interpreter

import (
	"github.com/DomBlack/lox/glox/pkg/ast"
	"github.com/DomBlack/lox/glox/pkg/errs"
)

var _ ast.StmtVisitor[any] = (*interpreter)(nil)

func (i *interpreter) execute(stmt ast.Stmt) {
	ast.AcceptStmt[any](stmt, i)
}

func (i *interpreter) executeBlock(stmts []ast.Stmt, environment *Environment) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *interpreter) VisitBlockStmt(stmt *ast.BlockStmt) any {
	i.executeBlock(stmt.Statements, i.environment.Scope())
	return nil
}

func (i *interpreter) VisitExpressionStmt(v *ast.ExpressionStmt) any {
	i.evaluate(v.Expression)
	return nil
}

func (i *interpreter) VisitIfStmt(v *ast.IfStmt) any {
	if isTruthy(i.evaluate(v.Condition)) {
		i.execute(v.ThenBranch)
	} else if v.ElseBranch != nil {
		i.execute(v.ElseBranch)
	}
	return nil
}

func (i *interpreter) VisitClassStmt(v *ast.ClassStmt) any {
	var superclass *Class
	if v.Superclass != nil {
		super, ok := i.evaluate(v.Superclass).(*Class)
		if !ok {
			panic(&errs.RuntimeError{v.Superclass.Name, "Superclass must be a class."})
		}
		superclass = super
	}

	i.environment.Define(v.Name.Lexeme, nil)

	if v.Superclass != nil {
		i.environment = i.environment.Scope()
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*Function)
	for _, method := range v.Methods {
		function := &Function{method, i.environment, method.Name.Lexeme == "init"}
		methods[method.Name.Lexeme] = function
	}

	class := &Class{v.Name.Lexeme, superclass, methods}

	if superclass != nil {
		i.environment = i.environment.Enclosing
	}

	i.environment.Assign(v.Name, class)
	return nil
}

func (i *interpreter) VisitFunctionStmt(v *ast.FunctionStmt) any {
	function := &Function{v, i.environment, false}
	i.environment.Define(v.Name.Lexeme, function)
	return nil
}

func (i *interpreter) VisitPrintStmt(v *ast.PrintStmt) any {
	value := i.evaluate(v.Expression)
	println(stringify(value))
	return nil
}

func (i *interpreter) VisitReturnStmt(v *ast.ReturnStmt) any {
	var value any
	if v.Value != nil {
		value = i.evaluate(v.Value)
	}

	panic(&Return{value})
}

func (i *interpreter) VisitWhileStmt(v *ast.WhileStmt) any {
	for isTruthy(i.evaluate(v.Condition)) {
		i.execute(v.Body)
	}
	return nil
}

func (i *interpreter) VisitVarStmt(v *ast.VarStmt) any {
	var value any
	if v.Initializer != nil {
		value = i.evaluate(v.Initializer)
	}

	i.environment.Define(v.Name.Lexeme, value)
	return nil
}
