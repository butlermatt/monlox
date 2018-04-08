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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.Boolean:
		return nativeBooltoObject(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return Null
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, node.Token.Line)
	}

	return Null
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result.Type() {
		case object.RETURN:
			return result.(*object.ReturnValue).Value
		case object.ERROR:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
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

func evalPrefixExpression(prefix *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(prefix.Right, env)

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

func evalInfixExpression(infix *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(infix.Left, env)
	right := Eval(infix.Right, env)

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return evalNumberInfixExpression(infix, left, right)
	}

	if left.Type() == object.STRING && right.Type() == object.STRING {
		return evalStringInfixExpression(infix, left, right)
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

func evalStringInfixExpression(infix *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch infix.Operator {
	case "==":
		return nativeBooltoObject(leftVal == rightVal)
	case "!=":
		return nativeBooltoObject(leftVal != rightVal)
	case "+":
		return &object.String{Value: leftVal + rightVal}
	}

	return newError(infix.Token.Line, "unknown operator: %s %s %s", left.Type(), infix.Operator, right.Type())
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.Token.Line, "identifier not found: "+node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		ev := Eval(e, env)
		if isError(ev) {
			return []object.Object{ev}
		}
		result = append(result, ev)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object, line int) object.Object {

	switch fn.Type() {
	case object.FUNCTION:
		function := fn.(*object.Function)
		exEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, exEnv)
		return unwrapReturnValue(evaluated)
	case object.BUILTIN:
		return fn.(*object.Builtin).Fn(line, args...)
	}

	return newError(line, "not a function: %s", fn.Type())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclodedEnvironment(fn.Env)

	for i, p := range fn.Parameters {
		env.Set(p.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if obj.Type() == object.RETURN {
		return obj.(*object.ReturnValue).Value
	}

	return obj
}
