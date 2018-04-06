package object

import (
	"bytes"
	"fmt"
	"github.com/butlermatt/monlox/ast"
	"strings"
)

type Type int

func (t Type) String() string {
	switch t {
	case NULL:
		return "NULL"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case STRING:
		return "STRING"
	case RETURN:
		return "RETURN"
	case FUNCTION:
		return "FUNCTION"
	case ERROR:
		return "ERROR"
	}

	return ""
}

const (
	NULL Type = iota
	NUMBER
	BOOLEAN
	STRING
	RETURN
	FUNCTION
	ERROR
)

type Object interface {
	Type() Type
	Inspect() string
}

type Number struct {
	Value float32
}

func (n *Number) Type() Type      { return NUMBER }
func (n *Number) Inspect() string { return fmt.Sprintf("%v", n.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BOOLEAN }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%v", b.Value) }

type Null struct{}

func (n *Null) Type() Type      { return NULL }
func (n *Null) Inspect() string { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type      { return RETURN }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Error struct {
	Message string
	Line    int
}

func (e *Error) Type() Type      { return ERROR }
func (e *Error) Inspect() string { return fmt.Sprintf("ERROR line %d: %s", e.Line, e.Message) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type { return FUNCTION }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type      { return STRING }
func (s *String) Inspect() string { return s.Value }
