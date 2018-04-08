package evaluator

import "github.com/butlermatt/monlox/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "wrong number of arguments. expected=%d, got=%d", 1, len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Number{Value: float32(len(arg.Value))}
			}

			return newError(line, "argument to `len` not supported. got=%s", args[0].Type())
		},
	},
}
