package evaluator

import (
	"fmt"
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
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result.Type() {
		case object.RETURN:
			return result.(*object.ReturnValue).Value
		case object.ERROR:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
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

func evalPrefixExpression(prefix *ast.PrefixExpression) object.Object {
	right := Eval(prefix.Right)

	switch prefix.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(prefix, right)
	default:
		return newError(prefix.Token.Line, "unknown operator: %s%s", prefix.Operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if right == False || right == Null {
		return True
	}

	return False
}

func evalMinusPrefixOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	if right.Type() != object.NUMBER {
		return newError(node.Token.Line, "unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalInfixExpression(infix *ast.InfixExpression) object.Object {
	left := Eval(infix.Left)
	right := Eval(infix.Right)

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return evalNumberInfixExpression(infix, left, right)
	}

	if left.Type() != right.Type() {
		return newError(infix.Token.Line, "type mismatch: %s %s %s", left.Type(), infix.Operator, right.Type())
	}

	switch infix.Operator {
	case "==":
		return nativeBooltoObject(left == right)
	case "!=":
		return nativeBooltoObject(left != right)
	case "or":
		return nativeBooltoObject((left == True) || (right == True))
	case "and":
		return nativeBooltoObject((left == True) && (right == True))
	}

	return newError(infix.Token.Line, "unknown operator: %s %s %s", left.Type(), infix.Operator, right.Type())
}

func evalNumberInfixExpression(infix *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	var result float32
	switch infix.Operator {
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
		return newError(infix.Token.Line, "unknown operator: %s %s %s", left.Type(), infix.Operator, right.Type())
	}

	return &object.Number{Value: result}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

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

func newError(line int, format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	return &object.Error{Line: line, Message: fmt.Sprintf("on line %d: %s", line, msg)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}

	return false
}
