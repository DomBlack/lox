package interpreter

import (
	"github.com/DomBlack/lox/glox/pkg/ast"
	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

type resolver struct {
	interpreter  *interpreter
	scopes       []map[string]bool
	currentFunc  FunctionType
	currentClass ClassType
}

var _ ast.StmtVisitor[any] = (*resolver)(nil)
var _ ast.ExprVisitor[any] = (*resolver)(nil)

func (r *resolver) resolve(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *resolver) resolveStmt(stmt ast.Stmt) {
	ast.AcceptStmt[any](stmt, r)
}

func (r *resolver) resolveExpr(expr ast.Expr) {
	ast.AcceptExpr[any](expr, r)
}

func (r *resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *resolver) declare(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		errs.ErrorAtToken(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false
}

func (r *resolver) define(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *resolver) resolveLocal(expr ast.Expr, name *token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *resolver) resolveFunction(fn *ast.FunctionStmt, funcType FunctionType) {
	enclosingFunc := r.currentFunc
	r.currentFunc = funcType

	r.beginScope()
	for _, param := range fn.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(fn.Body)
	r.endScope()

	r.currentFunc = enclosingFunc
}

func (r *resolver) VisitBlockStmt(v *ast.BlockStmt) any {
	r.beginScope()
	r.resolve(v.Statements)
	r.endScope()
	return nil
}

func (r *resolver) VisitExpressionStmt(v *ast.ExpressionStmt) any {
	r.resolveExpr(v.Expression)
	return nil
}

func (r *resolver) VisitClassStmt(v *ast.ClassStmt) any {
	enclosingClass := r.currentClass
	r.currentClass = CTClass

	r.declare(v.Name)
	r.define(v.Name)

	if v.Superclass != nil && v.Name.Lexeme == v.Superclass.Name.Lexeme {
		errs.ErrorAtToken(v.Superclass.Name, "A class cannot inherit from itself.")
	}

	if v.Superclass != nil {
		r.currentClass = CTSubclass
		r.resolveExpr(v.Superclass)
	}

	if v.Superclass != nil {
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range v.Methods {
		declaration := FTMethod
		if method.Name.Lexeme == "init" {
			declaration = FTInitializer
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()

	if v.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
	return nil
}

func (r *resolver) VisitFunctionStmt(v *ast.FunctionStmt) any {
	r.declare(v.Name)
	r.define(v.Name)
	r.resolveFunction(v, FTFunction)
	return nil
}

func (r *resolver) VisitIfStmt(v *ast.IfStmt) any {
	r.resolveExpr(v.Condition)
	r.resolveStmt(v.ThenBranch)
	if v.ElseBranch != nil {
		r.resolveStmt(v.ElseBranch)
	}
	return nil
}

func (r *resolver) VisitPrintStmt(v *ast.PrintStmt) any {
	r.resolveExpr(v.Expression)
	return nil
}

func (r *resolver) VisitReturnStmt(v *ast.ReturnStmt) any {
	if r.currentFunc == FTNone {
		errs.ErrorAtToken(v.Keyword, "Cannot return from top-level code.")
	}

	if v.Value != nil {
		if r.currentFunc == FTInitializer {
			errs.ErrorAtToken(v.Keyword, "Cannot return a value from an initializer.")
		}
		r.resolveExpr(v.Value)
	}
	return nil
}

func (r *resolver) VisitVarStmt(v *ast.VarStmt) any {
	r.declare(v.Name)
	if v.Initializer != nil {
		r.resolveExpr(v.Initializer)
	}
	r.define(v.Name)
	return nil
}

func (r *resolver) VisitWhileStmt(v *ast.WhileStmt) any {
	r.resolveExpr(v.Condition)
	r.resolveStmt(v.Body)
	return nil
}

func (r *resolver) VisitAssignExpr(v *ast.AssignExpr) any {
	r.resolveExpr(v.Value)
	r.resolveLocal(v, v.Name)
	return nil
}

func (r *resolver) VisitBinaryExpr(v *ast.BinaryExpr) any {
	r.resolveExpr(v.Left)
	r.resolveExpr(v.Right)
	return nil
}

func (r *resolver) VisitCallExpr(v *ast.CallExpr) any {
	r.resolveExpr(v.Callee)
	for _, arg := range v.Arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *resolver) VisitGetExpr(v *ast.GetExpr) any {
	r.resolveExpr(v.Object)
	return nil
}

func (r *resolver) VisitGroupingExpr(v *ast.GroupingExpr) any {
	r.resolveExpr(v.Expression)
	return nil
}

func (r *resolver) VisitLogicalExpr(v *ast.LogicalExpr) any {
	r.resolveExpr(v.Left)
	r.resolveExpr(v.Right)
	return nil
}

func (r *resolver) VisitSetExpr(v *ast.SetExpr) any {
	r.resolveExpr(v.Value)
	r.resolveExpr(v.Object)
	return nil
}

func (r *resolver) VisitSuperExpr(v *ast.SuperExpr) any {
	if r.currentClass == CTNone {
		errs.ErrorAtToken(v.Keyword, "Cannot use 'super' outside of a class.")
	} else if r.currentClass != CTSubclass {
		errs.ErrorAtToken(v.Keyword, "Cannot use 'super' in a class with no superclass.")
	}

	r.resolveLocal(v, v.Keyword)
	return nil
}

func (r *resolver) VisitThisExpr(v *ast.ThisExpr) any {
	if r.currentClass == CTNone {
		errs.ErrorAtToken(v.Keyword, "Cannot use 'this' outside of a class.")
		return nil
	}

	r.resolveLocal(v, v.Keyword)
	return nil
}

func (r *resolver) VisitLiteralExpr(v *ast.LiteralExpr) any {
	return nil
}

func (r *resolver) VisitUnaryExpr(v *ast.UnaryExpr) any {
	r.resolveExpr(v.Right)
	return nil
}

func (r *resolver) VisitVariableExpr(v *ast.VariableExpr) any {
	if len(r.scopes) != 0 {
		if defined, found := r.scopes[len(r.scopes)-1][v.Name.Lexeme]; !defined && found {
			errs.ErrorAtToken(v.Name, "Cannot read local variable in its own initializer.")
		}
	}

	r.resolveLocal(v, v.Name)
	return nil
}
