package lexer

import "github.com/butlermatt/monlox/token"

// Lexer iterates through the provided program to generate tokens.
type Lexer struct {
	input        string
	position     int  // Current position in the input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current character
	line         int  // current line
}

func newToken(tokenType token.TokenType, ch byte, line int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line}
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

// New returns a new Lexer populated with the specified input program.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isAlphaNumeric(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, bool) {
	pos := l.position + 1
	ok := true
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == 0 {
			ok = false
			break
		}
	}
	return l.input[pos:l.position], ok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line += 1
		}
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken steps through the input to generate the next token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch), l.line)
	case '(':
		tok = token.New(token.LPAREN, string(l.ch), l.line)
	case ')':
		tok = token.New(token.RPAREN, string(l.ch), l.line)
	case '{':
		tok = token.New(token.LBRACE, string(l.ch), l.line)
	case '}':
		tok = token.New(token.RBRACE, string(l.ch), l.line)
	case ',':
		tok = token.New(token.COMMA, string(l.ch), l.line)
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.EQ_EQ, literal, l.line)
		} else {
			tok = token.New(token.EQ, string(l.ch), l.line)
		}
	case '+':
		tok = token.New(token.PLUS, string(l.ch), l.line)
	case '-':
		tok = token.New(token.MINUS, string(l.ch), l.line)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.NOT_EQ, literal, l.line)
		} else {
			tok = token.New(token.BANG, string(l.ch), l.line)
		}
	case '/':
		tok = token.New(token.SLASH, string(l.ch), l.line)
	case '*':
		tok = token.New(token.ASTERISK, string(l.ch), l.line)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.New(token.LT_EQ, lit, l.line)
		} else {
			tok = token.New(token.LT, string(l.ch), l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.New(token.GT_EQ, lit, l.line)
		} else {
			tok = token.New(token.GT, string(l.ch), l.line)
		}
	case '"':
		start := l.line
		if str, ok := l.readString(); ok {
			tok = token.New(token.STRING, str, start)
		} else {
			tok = token.New(token.UTSTRING, str, start)
		}
	case 0:
		tok = token.New(token.EOF, "", l.line)
	default:
		if isAlpha(l.ch) {
			lit := l.readIdentifier()
			tok = token.New(token.LookupIdent(lit), lit, l.line)
			return tok
		} else if isDigit(l.ch) {
			tok := token.New(token.NUM, l.readNumber(), l.line)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	}

	l.readChar()
	return tok
}
