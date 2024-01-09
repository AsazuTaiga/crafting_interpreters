package main

import (
	"github.com/AsazuTaiga/crafting_interpriters/go/ast"
	"github.com/AsazuTaiga/crafting_interpriters/go/logger"
	"github.com/AsazuTaiga/crafting_interpriters/go/parser"
	"github.com/AsazuTaiga/crafting_interpriters/go/scanner"
	"github.com/AsazuTaiga/crafting_interpriters/go/token"
)

type AstPrinter struct {}

func (p AstPrinter) Print(expr ast.Expr) string {
	return expr.Accept(p).(string)
}

func (p AstPrinter) VisitBinaryExpr(expr ast.BinaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) VisitGroupingExpr(expr ast.GroupingExpr) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expr ast.LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.(string)
}

func (p AstPrinter) VisitUnaryExpr(expr ast.UnaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	str := "(" + name
	for _, expr := range exprs {
		str += " "
		str += expr.Accept(p).(string)
	}
	str += ")"
	return str
}


func main() {
	expression := ast.BinaryExpr{
		Left: &ast.UnaryExpr{
			Operator: token.Token{
				Type: token.MINUS,
				Lexeme: "-",
				Literal: nil,
				Line: 1,
			},
			Right: &ast.LiteralExpr{
				Value: "123",
			},
		},
		Operator: token.Token{
			Type: token.STAR,
			Lexeme: "*",
			Literal: nil,
			Line: 1,
		},
		Right: &ast.GroupingExpr{
			Expression: &ast.LiteralExpr{
				Value: "45.67",
			},
		},
	}

	printer := AstPrinter{}
	println(printer.Print(&expression))

	s := scanner.NewScanner("1 + 2")
	tokens := s.ScanTokens(logger.NewLogger())
	p := parser.NewParser(tokens)
	exp := p.Parse()
	println(printer.Print(exp))
}