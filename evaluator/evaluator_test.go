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

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 >= 2) { 10 }", nil},
		{"if (1 >= 2) { 10 } else { 20 }", 20},
		{"if (1 == 2 or 1 <= 2) { 10 } else { 20 }", 10},
		{"if (true and 1 == 1) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		number, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float32(number))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float32
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "on line 1: type mismatch: NUMBER + BOOLEAN"},
		{"5 + true; 5;", "on line 1: type mismatch: NUMBER + BOOLEAN"},
		{"-true;", "on line 1: unknown operator: -BOOLEAN"},
		{"true + false;", "on line 1: unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "on line 1: unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "on line 1: unknown operator: BOOLEAN + BOOLEAN"},
		{"if (1 == true) { 10 }", "on line 1: type mismatch: NUMBER == BOOLEAN"},
		{`
if (10 > 1) {
   if (10 > 1) { 
     return true + false 
   } 
   return 1; 
}`, "on line 4: unknown operator: BOOLEAN + BOOLEAN"},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("test %d: wrong type returned. expected=*object.Error, got=%T (%+v)", i+1, evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("test %d: wrong error message. expected=%q, got=%q", i+1, tt.expected, errObj.Message)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != Null {
		t.Errorf("object is not expected type. expected=Null, got=%T (%+v)", obj, obj)
		return false
	}

	return true
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
