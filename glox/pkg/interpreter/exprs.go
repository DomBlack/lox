package interpreter

import (
	"fmt"

	"github.com/DomBlack/lox/glox/pkg/ast"
	"github.com/DomBlack/lox/glox/pkg/errs"
	"github.com/DomBlack/lox/glox/pkg/token"
)

var _ ast.ExprVisitor[any] = (*interpreter)(nil)

func (i *interpreter) evaluate(expr ast.Expr) any {
	return ast.AcceptExpr[any](expr, i)
}

func (i *interpreter) VisitAssignExpr(v *ast.AssignExpr) any {
	value := i.evaluate(v.Value)

	distance, ok := i.locals[v]
	if ok {
		i.environment.AssignAt(distance, v.Name, value)
	} else {
		globals.Assign(v.Name, value)
	}

	return value
}

func (i *interpreter) VisitBinaryExpr(v *ast.BinaryExpr) any {
	left := i.evaluate(v.Left)
	right := i.evaluate(v.Right)

	switch v.Operator.Type {
	case token.GREATER:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	case token.MINUS:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) - right.(float64)
	case token.SLASH:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) / right.(float64)
	case token.STAR:
		checkNumberOperands(v.Operator, left, right)
		return left.(float64) * right.(float64)
	case token.PLUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right
			}
		}

		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right
			}
		}

		panic(&errs.RuntimeError{Token: v.Operator, Msg: "Operands must be two numbers or two strings."})
	}

	return nil
}

func (i *interpreter) VisitGetExpr(v *ast.GetExpr) any {
	object := i.evaluate(v.Object)
	if object, ok := object.(*Instance); ok {
		return object.Get(v.Name)
	}

	panic(&errs.RuntimeError{Token: v.Name, Msg: "Only instances have properties."})
}

func (i *interpreter) VisitGroupingExpr(v *ast.GroupingExpr) any {
	return i.evaluate(v.Expression)
}

func (i *interpreter) VisitLogicalExpr(v *ast.LogicalExpr) any {
	left := i.evaluate(v.Left)

	if v.Operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(v.Right)
}

func (i *interpreter) VisitLiteralExpr(v *ast.LiteralExpr) any {
	return v.Value
}

func (i *interpreter) VisitCallExpr(v *ast.CallExpr) any {
	callee := i.evaluate(v.Callee)

	var args []any
	for _, arg := range v.Arguments {
		args = append(args, i.evaluate(arg))
	}

	if function, ok := callee.(Callable); ok {
		if len(args) != function.Arity() {
			panic(&errs.RuntimeError{Token: v.Paren, Msg: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(args))})
		}

		return function.Call(i, args)
	}

	panic(&errs.RuntimeError{Token: v.Paren, Msg: "Can only call functions and classes."})
}

func (i *interpreter) VisitSetExpr(v *ast.SetExpr) any {
	object := i.evaluate(v.Object)

	instance, ok := object.(*Instance)
	if !ok {
		panic(&errs.RuntimeError{Token: v.Name, Msg: "Only instances have fields."})
	}

	value := i.evaluate(v.Value)
	instance.Set(v.Name, value)
	return value
}

func (i *interpreter) VisitSuperExpr(v *ast.SuperExpr) any {
	distance := i.locals[v]
	superclass := i.environment.GetAt(distance, "super").(*Class)

	object := i.environment.GetAt(distance-1, "this").(*Instance)

	method := superclass.findMethod(v.Method.Lexeme)
	if method == nil {
		panic(&errs.RuntimeError{Token: v.Method, Msg: "Undefined method '" + v.Method.Lexeme + "'."})
	}

	return method.Bind(object)
}

func (i *interpreter) VisitThisExpr(v *ast.ThisExpr) any {
	return i.lookupVariable(v.Keyword, v)
}

func (i *interpreter) VisitUnaryExpr(v *ast.UnaryExpr) any {
	right := i.evaluate(v.Right)

	switch v.Operator.Type {
	case token.BANG:
		return !isTruthy(right)
	case token.MINUS:
		checkNumberOperands(v.Operator, right)
		return -right.(float64)
	}

	return nil
}

func (i *interpreter) VisitVariableExpr(v *ast.VariableExpr) any {
	return i.lookupVariable(v.Name, v)
}
