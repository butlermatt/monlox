package evaluator

import (
	"github.com/butlermatt/monlox/ast"
	"github.com/butlermatt/monlox/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.Boolean:
		return nativeBooltoObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBooltoObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return Null
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if right == False || right == Null {
		return True
	}

	return False
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER {
		return Null
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return evalNumberInfixExpression(operator, left, right)
	}

	return Null
}

func evalNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	var result float32
	switch operator {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		result = leftVal / rightVal
	case "<":
		return nativeBooltoObject(leftVal < rightVal)
	case ">":
		return nativeBooltoObject(leftVal > rightVal)
	case "<=":
		return nativeBooltoObject(leftVal <= rightVal)
	case ">=":
		return nativeBooltoObject(leftVal >= rightVal)
	case "==":
		return nativeBooltoObject(leftVal == rightVal)
	case "!=":
		return nativeBooltoObject(leftVal != rightVal)
	default:
		return Null
	}

	return &object.Number{Value: result}
}
