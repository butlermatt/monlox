package ast

import (
	"bytes"

	"github.com/butlermatt/monlox/token"
	"strings"
)

// Node is a node within the AST tree.
type Node interface {
	// TokenLiteral returns the string literal of the token associated with this ast node.
	TokenLiteral() string
	String() string
}

// Statement represents an AST statement node.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an AST expression node.
type Expression interface {
	Node
	expressionNode()
}

// Program represents the statements comprising nodes of the AST tree.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the string literal of the token associated with this ast node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a string representation of the program.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement is an AST node representing a variable assignment
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the string literal of the token associated with this ast node.
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// String returns a string representation of the Let statement.
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteByte(';')

	return out.String()
}

// Identifier represents an variable identifier
type Identifier struct {
	Token token.Token // The token.IDENT token.
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the string literal of the token associated with this ast node.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns a string representation of the identifier.
func (i *Identifier) String() string { return i.Value }

// ReturnStatement is an AST node representing just the return statement and the associated expression.
type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the string literal of the token associated with this ast node.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a string representation of the Return statement
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteByte(';')

	return out.String()
}

// ExpressionStatement is a AST node representing a statement that consists of a single expression.
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the string literal of the token associated with this ast node.
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a string representation of an Expression statement.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// NumberLiteral is an AST node representing a number literal. Stored as a float32.
type NumberLiteral struct {
	Token token.Token
	Value float32
}

func (nl *NumberLiteral) expressionNode() {}

// TokenLiteral returns a string representation of the token associated with this node.
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }

// String returns a string representation of the Number Literal.
func (nl *NumberLiteral) String() string { return nl.Token.Literal }

// Boolean is an AST node representing boolean literals.
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral returns a string representation of the token associated with this node.
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns a string representation of the boolean literal.
func (b *Boolean) String() string { return b.Token.Literal }

// PrefixExpression is an AST node representing a prefix expression such as -5 or !x
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the string literal of the associated token.
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a string representation of the prefix expression.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteByte(')')

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {}

// TokenLiteral returns the string representation of this token.
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

// String return a string representation of this expression.
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteByte(')')

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // The { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral returns the string representation of this token.
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns the string representation of this expression.
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteByte(' ')
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token.
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the string representation of this token.
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns the string representation of this expression
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteByte('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral returns the string representation of this token.
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns a string representation of this call expression.
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteByte('(')
	out.WriteString(strings.Join(args, ", "))
	out.WriteByte(')')

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral returns a string representation of this token.
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String returns a string representation of this string literal
func (sl *StringLiteral) String() string { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // The '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteByte('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteByte(']')

	return out.String()
}

type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(ie.Left.String())
	out.WriteByte('[')
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // The '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteByte('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteByte('}')

	return out.String()
}
