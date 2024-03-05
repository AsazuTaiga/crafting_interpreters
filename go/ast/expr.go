package ast

import (
	"github.com/AsazuTaiga/crafting_interpriters/go/token"
)

type Expr interface {
	Accept(visitor Visitor) interface{}
}

type Visitor interface {
	VisitBinaryExpr(expr BinaryExpr) interface{}
	VisitGroupingExpr(expr GroupingExpr) interface{}
	VisitLiteralExpr(expr LiteralExpr) interface{}
	VisitUnaryExpr(expr UnaryExpr) interface{}
}

type BinaryExpr struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinaryExpr(left Expr, operator token.Token, right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left: left,
		Operator: operator,
		Right: right,
	}
}

func (expr *BinaryExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpr(*expr)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{
		Expression: expression,
	}
}

func (expr *GroupingExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpr(*expr)
}

type LiteralExpr struct {
	Value interface{}
}

func NewLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{
		Value: value,
	}
}

func (expr *LiteralExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpr(*expr)
}

type UnaryExpr struct {
	Operator token.Token
	Right Expr
}

func NewUnaryExpr(operator token.Token, right Expr) *UnaryExpr {
	return &UnaryExpr{
		Operator: operator,
		Right: right,
	}
}

func (expr *UnaryExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(*expr)
}


type VariableExpr struct {
	Name token.Token
}

func NewVariableExpr(name token.Token) *VariableExpr {
	return &VariableExpr{
		Name: name,
	}
}

func (expr *VariableExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariableExpr(*expr)
}
