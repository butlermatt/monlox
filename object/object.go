package object

import "fmt"

type Type int

func (t Type) String() string {
	switch t {
	case NULL:
		return "NULL"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case RETURN:
		return "RETURN"
	case ERROR:
		return "ERROR"
	}

	return ""
}

const (
	NULL Type = iota
	NUMBER
	BOOLEAN
	RETURN
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
