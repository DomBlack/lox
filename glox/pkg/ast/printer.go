package ast

import (
	"fmt"
	"strings"
)

func Print(expr Expr) string {
	return AcceptExpr[string](expr, &printer{})
}

type printer struct {
}

var _ ExprVisitor[string] = (*printer)(nil)

func (p *printer) Print(expr Expr) string {
	return AcceptExpr[string](expr, p)
}

func (p *printer) VisitAssignExpr(v *AssignExpr) string {
	return fmt.Sprintf("%s = %s", v.Name.Lexeme, p.Print(v.Value))
}

func (p *printer) VisitBinaryExpr(v *BinaryExpr) string {
	return parenthesize(v.Operator.Lexeme, v.Left, v.Right)
}

func (p *printer) VisitGroupingExpr(v *GroupingExpr) string {
	return parenthesize("group", v.Expression)
}

func (p *printer) VisitLiteralExpr(v *LiteralExpr) string {
	if v.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", v.Value)
}

func (p *printer) VisitLogicalExpr(v *LogicalExpr) string {
	return parenthesize(v.Operator.Lexeme, v.Left, v.Right)
}

func (p *printer) VisitCallExpr(v *CallExpr) string {
	var builder strings.Builder

	builder.WriteString(AcceptExpr[string](v.Callee, p))
	builder.WriteString("(")
	for i, arg := range v.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(AcceptExpr[string](arg, p))
	}
	builder.WriteString(")")

	return builder.String()
}

func (p *printer) VisitGetExpr(v *GetExpr) string {
	return fmt.Sprintf("%s.%s", p.Print(v.Object), v.Name.Lexeme)
}

func (p *printer) VisitSetExpr(v *SetExpr) string {
	return fmt.Sprintf("%s.%s = %s", p.Print(v.Object), v.Name.Lexeme, p.Print(v.Value))
}

func (p *printer) VisitSuperExpr(v *SuperExpr) string {
	return fmt.Sprintf("super.%s", v.Method.Lexeme)
}

func (p *printer) VisitThisExpr(_ *ThisExpr) string {
	return "this"
}

func (p *printer) VisitUnaryExpr(v *UnaryExpr) string {
	return parenthesize(v.Operator.Lexeme, v.Right)
}

func (p *printer) VisitVariableExpr(v *VariableExpr) string {
	return v.Name.Lexeme
}

func parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(AcceptExpr[string](expr, &printer{}))
	}
	builder.WriteString(")")

	return builder.String()
}
