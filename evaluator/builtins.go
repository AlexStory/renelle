// evaluator/builtins.go

package evaluator

import (
	"fmt"
	"renelle/constants"
	"renelle/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(ctx.Line, ctx.Column, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"head": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(ctx.Line, ctx.Column, "argument to `head` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return constants.NIL
		},
	},

	"tail": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(ctx.Line, ctx.Column, "argument to `tail` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return constants.NIL
		},
	},
	"last": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(ctx.Line, ctx.Column, "argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return constants.NIL
		},
	},
	"push": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(ctx.Line, ctx.Column, "argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"fst": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.TUPLE_OBJ {
				return newError(ctx.Line, ctx.Column, "argument to `fst` must be TUPLE, got %s", args[0].Type())
			}

			tuple := args[0].(*object.Tuple)
			if len(tuple.Elements) > 0 {
				return tuple.Elements[0]
			}

			return constants.NIL
		},
	},
	"snd": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.TUPLE_OBJ {
				return newError(ctx.Line, ctx.Column, "argument to `snd` must be TUPLE, got %s", args[0].Type())
			}

			tuple := args[0].(*object.Tuple)
			if len(tuple.Elements) > 1 {
				return tuple.Elements[1]
			}

			return constants.NIL
		},
	},
	"print": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}

			fmt.Println(args[0].Inspect())
			return constants.OK
		},
	},
	"os_args": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=0", len(args))
			}

			arguments := &object.Array{}

			for _, arg := range (*ctx.MetaData)["args"].([]string) {
				arguments.Elements = append(arguments.Elements, &object.String{Value: arg})
			}

			return arguments
		},
	},
	"type": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=1", len(args))
			}

			return &object.String{Value: string(args[0].Type())}
		},
	},
}
