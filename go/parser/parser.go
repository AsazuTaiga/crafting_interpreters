package parser

import (
	"github.com/AsazuTaiga/crafting_interpriters/go/ast"
	"github.com/AsazuTaiga/crafting_interpriters/go/logger"
	"github.com/AsazuTaiga/crafting_interpriters/go/stmt"

	"github.com/AsazuTaiga/crafting_interpriters/go/token"
)

type Parser struct {
	tokens []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens: tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []stmt.Stmt {
	var statements []stmt.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() stmt.Stmt {
	if(p.match(token.VAR)) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() stmt.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")
	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variavle declaration.")
	return &stmt.Var{
		Name: *name,
		Initializer: initializer,
	}
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) statement() stmt.Stmt {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() stmt.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &stmt.Print{
		Expression: value,
	}
}

func (p *Parser) expressionStatement() stmt.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &stmt.Expression{
		Expression: value,
	}
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinaryExpr(expr, *operator, right)
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinaryExpr(
			expr,
			*operator,
			right,
		)
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		token := p.previous()
		right := p.factor()
		expr = ast.NewBinaryExpr(
			expr,
			*token,
			right,
		)
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		token := p.previous()
		right := p.unary()
		expr = ast.NewBinaryExpr(
			expr,
			*token,
			right,
	)
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnaryExpr(
			*operator,
			right,
		)
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpr(false)
	}
	if p.match(token.TRUE) {
		return ast.NewLiteralExpr(true)
	}
	if p.match(token.NIL) {
		return ast.NewLiteralExpr(nil)
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteralExpr(p.previous().Literal)
	}

	if p.match(token.IDENTIFIER) {
		return ast.NewVariableExpr(*p.previous())
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.NewGroupingExpr(expr)
	}

	p.error(p.peek(), "Expect expression.")

	return nil
}

func (p *Parser) consume(t token.TokenType, message string) *token.Token {
	if p.check(t) {
		return p.advance()
	}

	p.error(p.peek(), message)
	return nil
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.TokenType) bool {
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

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token *token.Token, message string) {
	logger := logger.NewLogger()
	logger.ErrorReport(token.Line, message)
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