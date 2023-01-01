// Code generated by generate_ast.go; DO NOT EDIT.

package ast

import (
  "github.com/DomBlack/lox/glox/pkg/token"
)

type Expr interface {
  _expr() // unexported interface
}

type ExprVisitor[R any] interface {
  VisitAssignExpr(v *AssignExpr) R
  VisitBinaryExpr(v *BinaryExpr) R
  VisitCallExpr(v *CallExpr) R
  VisitGetExpr(v *GetExpr) R
  VisitGroupingExpr(v *GroupingExpr) R
  VisitLogicalExpr(v *LogicalExpr) R
  VisitLiteralExpr(v *LiteralExpr) R
  VisitSetExpr(v *SetExpr) R
  VisitSuperExpr(v *SuperExpr) R
  VisitThisExpr(v *ThisExpr) R
  VisitUnaryExpr(v *UnaryExpr) R
  VisitVariableExpr(v *VariableExpr) R
}

func AcceptExpr[R any](e Expr, v ExprVisitor[R]) R {
  switch e := e.(type) {
  case *AssignExpr:
    return v.VisitAssignExpr(e)
  case *BinaryExpr:
    return v.VisitBinaryExpr(e)
  case *CallExpr:
    return v.VisitCallExpr(e)
  case *GetExpr:
    return v.VisitGetExpr(e)
  case *GroupingExpr:
    return v.VisitGroupingExpr(e)
  case *LogicalExpr:
    return v.VisitLogicalExpr(e)
  case *LiteralExpr:
    return v.VisitLiteralExpr(e)
  case *SetExpr:
    return v.VisitSetExpr(e)
  case *SuperExpr:
    return v.VisitSuperExpr(e)
  case *ThisExpr:
    return v.VisitThisExpr(e)
  case *UnaryExpr:
    return v.VisitUnaryExpr(e)
  case *VariableExpr:
    return v.VisitVariableExpr(e)
  default:
    panic("Unknown type")
  }
}

type AssignExpr struct {
  Name *token.Token
  Value Expr
}
var _ Expr = (*AssignExpr)(nil)

func (e *AssignExpr) _expr() {}

type BinaryExpr struct {
  Left Expr
  Operator *token.Token
  Right Expr
}
var _ Expr = (*BinaryExpr)(nil)

func (e *BinaryExpr) _expr() {}

type CallExpr struct {
  Callee Expr
  Paren *token.Token
  Arguments []Expr
}
var _ Expr = (*CallExpr)(nil)

func (e *CallExpr) _expr() {}

type GetExpr struct {
  Object Expr
  Name *token.Token
}
var _ Expr = (*GetExpr)(nil)

func (e *GetExpr) _expr() {}

type GroupingExpr struct {
  Expression Expr
}
var _ Expr = (*GroupingExpr)(nil)

func (e *GroupingExpr) _expr() {}

type LogicalExpr struct {
  Left Expr
  Operator *token.Token
  Right Expr
}
var _ Expr = (*LogicalExpr)(nil)

func (e *LogicalExpr) _expr() {}

type LiteralExpr struct {
  Value any
}
var _ Expr = (*LiteralExpr)(nil)

func (e *LiteralExpr) _expr() {}

type SetExpr struct {
  Object Expr
  Name *token.Token
  Value Expr
}
var _ Expr = (*SetExpr)(nil)

func (e *SetExpr) _expr() {}

type SuperExpr struct {
  Keyword *token.Token
  Method *token.Token
}
var _ Expr = (*SuperExpr)(nil)

func (e *SuperExpr) _expr() {}

type ThisExpr struct {
  Keyword *token.Token
}
var _ Expr = (*ThisExpr)(nil)

func (e *ThisExpr) _expr() {}

type UnaryExpr struct {
  Operator *token.Token
  Right Expr
}
var _ Expr = (*UnaryExpr)(nil)

func (e *UnaryExpr) _expr() {}

type VariableExpr struct {
  Name *token.Token
}
var _ Expr = (*VariableExpr)(nil)

func (e *VariableExpr) _expr() {}
