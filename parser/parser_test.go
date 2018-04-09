package parser

import (
	"fmt"
	"testing"

	"github.com/butlermatt/monlox/ast"
	"github.com/butlermatt/monlox/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value interface{}
	}{
		{"let x = 5;", "x", float32(5)},
		{"let x = 5.43;", "x", float32(5.43)},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program contains incorrect number of Statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.ident) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.value) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %q. got=%q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %q. got=%q", name, letStmt.Name)
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"return 5;", float32(5)},
		{"return true", true},
		{"return 5.43;", float32(5.43)},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program contains incorrect number of Statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]
		val := stmt.(*ast.ReturnStatement).Value
		if !testLiteralExpression(t, val, tt.value) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	tests := []struct {
		input string
		value string
	}{
		{"foo;", "foo"},
		{"bar;", "bar"},
		{"f00;", "f00"},
		{"b4r;", "b4r"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		testIdentifier(t, stmt.Expression, tt.value)
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		value float32
	}{
		{"5;", 5},
		{"10;", 10},
		{"123.456;", 123.456},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have correct number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not right type. expected=ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		testNumberLiteral(t, stmt.Expression, tt.value)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not right type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
		}
		testBooleanLiteral(t, stmt.Expression, tt.value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		oper  string
		value interface{}
	}{
		{"!15;", "!", float32(15)},
		{"-5.2;", "-", float32(5.2)},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not have correct number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is wrong type. expected=ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not correct type. expected=ast.PrefixExpression, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.oper {
			t.Fatalf("exp.Operator is incorrect. expected=%q, got=%q", tt.oper, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testNumberLiteral(t *testing.T, nl ast.Expression, value float32) bool {
	num, ok := nl.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("nl not correct type. expected=*ast.NumberLiteral, got=%T", nl)
		return false
	}

	if num.Value != value {
		t.Errorf("num.Value is incorrect. expected=%v, got=%v", value, num.Value)
		return false
	}

	if num.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("num.TokenLiteral incorrect. expected=%v, got=%q", value, num.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	b, ok := bl.(*ast.Boolean)
	if !ok {
		t.Errorf("bl is not correct type. expected=*ast.Boolean, got=%T", bl)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value is incorrect. expected=%v, got=%v", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("b.TokenLiteral incorrect. expected=%v, got=%q", value, b.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 4;", 5, "+", 4},
		{"5 - 4;", 5, "-", 4},
		{"5 * 4;", 5, "*", 4},
		{"5 / 4;", 5, "/", 4},
		{"5 > 4;", 5, ">", 4},
		{"5 < 4;", 5, "<", 4},
		{"5.5 == 4.5;", float32(5.5), "==", float32(4.5)},
		{"5 != 4;", 5, "!=", 4},
		{"5 >= 4;", 5, ">=", 4},
		{"5 <= 4;", 5, "<=", 4},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contains wrong number of values. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is wrong type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 <= 4 != 3 >= 4", "((5 <= 4) != (3 >= 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2.2 / (5.5 + 5)", "(2.2 / (5.5 + 5))"},
		{"-(5.5 + 5.5)", "(-(5.5 + 5.5))"},
		{"!(true == true)", "(!(true == true))"},
		{"3 < 5 or 5 > 3", "((3 < 5) or (5 > 3))"},
		{"x <= 5 and 5 >= x", "((x <= 5) and (5 >= x))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statement is wrong type. expected=*ast.ExpressionStatement, got=%T (%+[1]v)", program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expression is wrong type. expected=*ast.ArrayLiteral, got=%T (%+[1]v)", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array contains wrong number of elements. expected=%d, got=%d", 3, len(array.Elements))
	}

	testNumberLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statement wrong type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("statement expression wrong type. expected=*ast.IndexExpression, got=%T (%+[1]v)", stmt.Expression)
	}

	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp wrong type. expected=*ast.Identifier, got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value incorrect. expected=%s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral incorrect. expected=%s, got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case float32:
		return testNumberLiteral(t, exp, v)
	case float64:
	case int:
		return testNumberLiteral(t, exp, float32(v))
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is wrong type. expected=*ast.InfixExpression, got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator incorrect. expected=%s, got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain correct number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] wrong type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression wrong type. expected=*ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence does not contain correct number of statements. expected=%d, got=%d", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is wrong type. expected=*ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was wrong value. expected=<nil>, got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain correct number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] wrong type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression wrong type. expected=*ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence does not contain correct number of statements. expected=%d, got=%d", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is wrong type. expected=*ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Fatalf("alternative is wrong type. expected=*ast.BlockStatement, got=%T", nil)
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative does not contain correct number of statements. expected=%d, got=%d", 1, len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Statements[0] wrong type. expected=*ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain correct number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is wrong type. expcted=*ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	fun, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is wrong type. expected=*ast.FunctionLiteral, got=%T", stmt.Expression)
	}

	if len(fun.Parameters) != 2 {
		t.Fatalf("function.Parameters does not contain correct number of paramets. expected=%d, got=%d", 2, len(fun.Parameters))
	}

	testLiteralExpression(t, fun.Parameters[0], "x")
	testLiteralExpression(t, fun.Parameters[1], "y")

	if len(fun.Body.Statements) != 1 {
		t.Fatalf("function.Body contains incorrect number of statements. expected=%d, got=%d", 1, len(fun.Body.Statements))
	}

	bodyStmt, ok := fun.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body statement is incorrec type. expected=*ast.ExpressionStatement, got=%T", fun.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: `fn() {};`, expected: []string{}},
		{input: `fn(x) {};`, expected: []string{"x"}},
		{input: `fn(x, y, z) {};`, expected: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expected) {
			t.Errorf("wrong number of parameters. expected=%d, got=%d", len(tt.expected), len(function.Parameters))
		}

		for i, ident := range tt.expected {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements contains wrong number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is incorrect type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is incorrect type. expected=*ast.CallExpression, got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("exp does not contain correct number of Arguements. expected=%d, got=%d", 3, len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expression wrong type. expected=*ast.StringLiteral, got=%T, (%+[1]v)", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal value is wrong. expected=%q, got=%q", "hello world", literal.Value)
	}
}

func TestUnterminatedString(t *testing.T) {
	input := `"hello World`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	errs := p.Errors()
	if len(errs) != 1 {
		t.Fatalf("incorrect number of errors. expected=%d, got=%d", 1, len(errs))
	}

	expected := "on line 1: unterminated string"
	if errs[0] != expected {
		t.Errorf("unexpected error message. expected=%q, got=%q", expected, errs[0])
	}
}
