package interpreter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AsazuTaiga/crafting_interpriters/go/ast"
	"github.com/AsazuTaiga/crafting_interpriters/go/environment"
	"github.com/AsazuTaiga/crafting_interpriters/go/stmt"
	"github.com/AsazuTaiga/crafting_interpriters/go/token"
)

type Interpreter struct {
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: environment.NewEnvironment(),
	}
}

func(i *Interpreter) Interpret(statements []stmt.Stmt) interface{} {
	for _, statement := range statements {
		i.execute(statement)
	}
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr ast.LiteralExpr) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr ast.GroupingExpr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(expr ast.UnaryExpr) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.BANG:
		return !isTruthy(right);
	case token.MINUS:
		err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return -right.(float64)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr ast.VariableExpr) interface{} {
	name, ok := i.environment.Get(expr.Name.String())
	if !ok {
		err := errors.New(fmt.Sprintf("Error: %s", "Undefined variable '"+expr.Name.String()+"'."))
		fmt.Printf("%s\n", err)
		return nil
	}
	return name
}

func (i *Interpreter) VisitBinaryExpr(expr ast.BinaryExpr) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

  switch expr.Operator.Type {
	case token.GREATER:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) >= right.(float64)
	case token.LESS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) <= right.(float64)
	case token.MINUS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) - right.(float64)
	case token.PLUS:
		switch left.(type) {
		case float64:
			switch right.(type) {
			case float64:
				return left.(float64) + right.(float64)
			}
		case string:
			switch right.(type) {
			case string:
				return left.(string) + right.(string)
			}
		}
		err := errors.New(fmt.Sprintf("Error: %s", "Operands must be two numbers or two strings."))
		fmt.Printf("%s\n", err)
		return nil
	case token.SLASH:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) / right.(float64)
	case token.STAR:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		return left.(float64) * right.(float64)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

// func (i Interpreter) VisitBlock(block *stmt.Block) {
// 	i.executeBlock(block.Statements, newEnvironment(i.env))
// }

// func (i Interpreter) executeBlock(statements []stmt.Stmt, env environment) {
// 	previousEnv := i.env
// 	defer func() {
// 		i.env = previousEnv
// 	}()

// 	i.env = &env
// 	for _, statement := range statements {
// 		i.execute(statement)
// 	}
// }

func (i *Interpreter) execute(stmt stmt.Stmt) {
	stmt.Accept(i)
}

func(i *Interpreter) evaluate(expression ast.Expr) interface{} {
	return expression.Accept(i)
}

func (i *Interpreter) VisitStatementExpression(stmt stmt.Expression) {
	i.evaluate(stmt.Expression)
	return
}

func (i *Interpreter) VisitPrintStmt(stmt stmt.Print) {
	value := i.evaluate(stmt.Expression)
	fmt.Printf("%s\n", stringify(value))
	return
}

func  isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	switch object.(type) {
	case bool:
		return object.(bool)
	}
	return true
}

func  isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func  checkNumberOperand(operator token.Token, operand interface{}) error {
	switch operand.(type) {
	case float64:
		return nil
	}
	return errors.New(fmt.Sprintf("Error: %s", "Operand must be a number."))
}

func  checkNumberOperands(operator token.Token, left interface{}, right interface{}) error {
	switch left.(type) {
	case float64:
		switch right.(type) {
		case float64:
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Error: %s", "Operands must be numbers."))
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	switch object.(type) {
	case float64:
		text := fmt.Sprintf("%g", object)
		if strings.HasSuffix(text, ".0") {
			text = text[:len(text)-2]
		}
		return text
	}
	return object.(string)
}

