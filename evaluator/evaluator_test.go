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
		{"-5", -5},
		{"-10.45", -10.45},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5.5 * 2 + 10", 21},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * 2.5 * 10", 50},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"true != true", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 >= 2) == true", false},
		{"(1 >= 2) == false", true},
		{"true or true", true},
		{"true or false", true},
		{"false or true", true},
		{"true and true", true},
		{"true and false", false},
		{"false and false", false},
		{"false and true", false},
		//{"1 != true", false},
		//{"1 == true", true},
		//{"1 or 2", true},
		//{"1 and 2", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not expected type. expected=*object.Boolean, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%v, got=%v", expected, result.Value)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}
