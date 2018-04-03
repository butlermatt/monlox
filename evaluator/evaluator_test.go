package evaluator

import (
	"testing"

	"github.com/butlermatt/monlox/lexer"
	"github.com/butlermatt/monlox/object"
	"github.com/butlermatt/monlox/parser"
)

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float32
	}{
		{"5", 5},
		{"10.45", 10.45},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testNumberObject(t *testing.T, obj object.Object, expected float32) bool {
	result, ok := obj.(*object.Number)
	if !ok {
		t.Errorf("object is not expected type. expected=*object.Number, got=%T, (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%v, got=%v", expected, result.Value)
		return false
	}

	return true
}
