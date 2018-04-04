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
		return evalProgram(node)
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
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		if result.Type() == object.RETURN {
			return result.(*object.ReturnValue).Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == object.RETURN {
			return result
		}
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

	switch operator {
	case "==":
		return nativeBooltoObject(left == right)
	case "!=":
		return nativeBooltoObject(left != right)
	case "or":
		return nativeBooltoObject((left == True) || (right == True))
	case "and":
		return nativeBooltoObject((left == True) && (right == True))
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

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return Null
}

func isTruthy(obj object.Object) bool {
	if obj == Null || obj == False {
		return false
	}

	return true
}
