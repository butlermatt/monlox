package evaluator

import "github.com/butlermatt/monlox/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(line int, args ...object.Object) object.Object {
			if e := expectNArgs(line, 1, args); e != nil {
				return e
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Number{Value: float32(len(arg.Elements))}
			case *object.String:
				return &object.Number{Value: float32(len(arg.Value))}
			}

			return newError(line, "argument to `len` not supported. got=%s", args[0].Type())
		},
	},
	"first": {
		Fn: func(line int, args ...object.Object) object.Object {
			if e := expectNArgs(line, 1, args); e != nil {
				return e
			}

			if args[0].Type() != object.ARRAY {
				return newError(line, "argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arg := args[0].(*object.Array)
			if len(arg.Elements) > 0 {
				return arg.Elements[0]
			}

			return Null
		},
	},
	"last": {
		Fn: func(line int, args ...object.Object) object.Object {
			if e := expectNArgs(line, 1, args); e != nil {
				return e
			}

			if args[0].Type() != object.ARRAY {
				return newError(line, "argument to `last` must be ARRAY, got=%s", args[0].Type())
			}

			arg := args[0].(*object.Array)
			if length := len(arg.Elements); length > 0 {
				return arg.Elements[length-1]
			}

			return Null
		},
	},
	"rest": {
		Fn: func(line int, args ...object.Object) object.Object {
			if e := expectNArgs(line, 1, args); e != nil {
				return e
			}

			if args[0].Type() != object.ARRAY {
				return newError(line, "argument to `rest` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if length := len(arr.Elements); length > 0 {
				newEls := make([]object.Object, length-1, length-1)
				copy(newEls, arr.Elements[1:])
				return &object.Array{Elements: newEls}
			}

			return Null
		},
	},
	"push": {
		Fn: func(line int, args ...object.Object) object.Object {
			if e := expectNArgs(line, 2, args); e != nil {
				return e
			}

			if args[0].Type() != object.ARRAY {
				return newError(line, "first argument to `push` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newEls := make([]object.Object, length+1, length+1)
			copy(newEls, arr.Elements)
			newEls[length] = args[1]

			return &object.Array{Elements: newEls}
		},
	},
}

func expectNArgs(line, expect int, args []object.Object) *object.Error {
	if len(args) != expect {
		return newError(line, "wrong number of arguments. expected=%d, got=%d", expect, len(args))
	}

	return nil
}
