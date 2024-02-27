// evaluator/builtins.go

package evaluator

import "renelle/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError(line, col, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
