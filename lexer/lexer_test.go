package lexer

import (
	"testing"

	"github.com/butlermatt/monlox/token"
)

func TestNextToken(t *testing.T) {
	input := `let five5 = 5;
let ten = 10;
let float = 3.5;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

10 == 10 or 10 != 9;
10 <= 9 and 10 >= 9;
"foobar";
"foo bar";
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{token.LET, "let", 1},
		{token.IDENT, "five5", 1},
		{token.EQ, "=", 1},
		{token.NUM, "5", 1},
		{token.SEMICOLON, ";", 1},
		{token.LET, "let", 2},
		{token.IDENT, "ten", 2},
		{token.EQ, "=", 2},
		{token.NUM, "10", 2},
		{token.SEMICOLON, ";", 2},
		{token.LET, "let", 3},
		{token.IDENT, "float", 3},
		{token.EQ, "=", 3},
		{token.NUM, "3.5", 3},
		{token.SEMICOLON, ";", 3},

		{token.LET, "let", 5},
		{token.IDENT, "add", 5},
		{token.EQ, "=", 5},
		{token.FUNCTION, "fn", 5},
		{token.LPAREN, "(", 5},
		{token.IDENT, "x", 5},
		{token.COMMA, ",", 5},
		{token.IDENT, "y", 5},
		{token.RPAREN, ")", 5},
		{token.LBRACE, "{", 5},
		{token.IDENT, "x", 6},
		{token.PLUS, "+", 6},
		{token.IDENT, "y", 6},
		{token.SEMICOLON, ";", 6},
		{token.RBRACE, "}", 7},
		{token.SEMICOLON, ";", 7},

		{token.LET, "let", 9},
		{token.IDENT, "result", 9},
		{token.EQ, "=", 9},
		{token.IDENT, "add", 9},
		{token.LPAREN, "(", 9},
		{token.IDENT, "five", 9},
		{token.COMMA, ",", 9},
		{token.IDENT, "ten", 9},
		{token.RPAREN, ")", 9},
		{token.SEMICOLON, ";", 9},
		{token.BANG, "!", 10},
		{token.MINUS, "-", 10},
		{token.SLASH, "/", 10},
		{token.ASTERISK, "*", 10},
		{token.NUM, "5", 10},
		{token.SEMICOLON, ";", 10},
		{token.NUM, "5", 11},
		{token.LT, "<", 11},
		{token.NUM, "10", 11},
		{token.GT, ">", 11},
		{token.NUM, "5", 11},
		{token.SEMICOLON, ";", 11},

		{token.IF, "if", 13},
		{token.LPAREN, "(", 13},
		{token.NUM, "5", 13},
		{token.LT, "<", 13},
		{token.NUM, "10", 13},
		{token.RPAREN, ")", 13},
		{token.LBRACE, "{", 13},
		{token.RETURN, "return", 14},
		{token.TRUE, "true", 14},
		{token.SEMICOLON, ";", 14},
		{token.RBRACE, "}", 15},
		{token.ELSE, "else", 15},
		{token.LBRACE, "{", 15},
		{token.RETURN, "return", 16},
		{token.FALSE, "false", 16},
		{token.SEMICOLON, ";", 16},
		{token.RBRACE, "}", 17},

		{token.NUM, "10", 19},
		{token.EQ_EQ, "==", 19},
		{token.NUM, "10", 19},
		{token.OR, "or", 19},
		{token.NUM, "10", 19},
		{token.NOT_EQ, "!=", 19},
		{token.NUM, "9", 19},
		{token.SEMICOLON, ";", 19},
		{token.NUM, "10", 20},
		{token.LT_EQ, "<=", 20},
		{token.NUM, "9", 20},
		{token.AND, "and", 20},
		{token.NUM, "10", 20},
		{token.GT_EQ, ">=", 20},
		{token.NUM, "9", 20},
		{token.SEMICOLON, ";", 20},
		{token.STRING, "foobar", 21},
		{token.SEMICOLON, ";", 21},
		{token.STRING, "foo bar", 22},
		{token.SEMICOLON, ";", 22},

		{token.EOF, "", 23},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
		}
	}
}
