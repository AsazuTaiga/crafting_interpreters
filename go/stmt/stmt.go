package stmt

import (
	"github.com/AsazuTaiga/crafting_interpriters/go/ast"
	"github.com/AsazuTaiga/crafting_interpriters/go/token"
)


type Stmt interface {
	Accept(visitor Visitor)
}

type Expression struct {
	Expression ast.Expr
}

func (expression *Expression) Accept(visitor Visitor) {
	visitor.VisitStatementExpression(expression)
}

type Block struct {
	Statements []Stmt
}

func (block *Block) Accept(visitor Visitor) {
	visitor.VisitBlock(block)
}

type Conditional struct {
	Condition     ast.Expr
	ThenStatement Stmt
	ElseStatement Stmt
}

func (conditional *Conditional) Accept(visitor Visitor) {
	visitor.VisitStatementConditional(conditional)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (function *Function) Accept(visitor Visitor) {
	visitor.VisitStatementFunction(function)
}

type Print struct {
	Expression ast.Expr
}

func (p *Print) Accept(visitor Visitor)  {
	visitor.VisitStatementPrint(p)
}

type While struct {
	Condition ast.Expr
	Body      Stmt
}

func (while *While) Accept(visitor Visitor) {
	visitor.VisitStatementWhile(while)
}

type Var struct {
	Name        token.Token
	Initializer ast.Expr
}

func (v *Var) Accept(visitor Visitor) {
	visitor.VisitStatementVar(v)
}

type Visitor interface {
	VisitStatementExpression(expression *Expression)
	VisitStatementPrint(p *Print)
	VisitStatementVar(v *Var)
	VisitBlock(block *Block)
	VisitStatementConditional(conditional *Conditional)
	VisitStatementWhile(while *While)
	VisitStatementFunction(function *Function)
}