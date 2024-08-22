// evaluator/builtins.go

package evaluator

import (
	"fmt"
	"renelle/constants"
	"renelle/object"
)

func reduceWhile(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 3 {
		return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=2 or 3", len(args))
	}

	list, ok := args[0].(*object.Array)
	if !ok {
		return newError(ctx.Line, ctx.Column, "first argument to `reduce_while` must be ARRAY, got %s", args[0].Type())
	}

	var initial object.Object
	var fn *object.Function
	var startIndex int

	if len(args) == 2 {
		if len(list.Elements) == 0 {
			return newError(ctx.Line, ctx.Column, "cannot reduce empty array without initial value")
		}
		initial = list.Elements[0]
		fn, ok = args[1].(*object.Function)
		if !ok {
			return newError(ctx.Line, ctx.Column, "second argument to `reduce_while` must be FUNCTION, got %s", args[1].Type())
		}
		startIndex = 1
	} else {
		initial = args[1]
		fn, ok = args[2].(*object.Function)
		if !ok {
			return newError(ctx.Line, ctx.Column, "third argument to `reduce_while` must be FUNCTION, got %s", args[2].Type())
		}
		startIndex = 0
	}

	accumulator := initial
	for i := startIndex; i < len(list.Elements); i++ {
		elem := list.Elements[i]
		result := applyFunction(fn, []object.Object{accumulator, elem}, ctx)
		if result.Type() == object.ERROR_OBJ {
			return result
		}

		tuple, ok := result.(*object.Tuple)
		if !ok || len(tuple.Elements) != 2 {
			return newError(ctx.Line, ctx.Column, "function must return a tuple of (:cont, acc) or (:halt, acc)")
		}

		action, ok := tuple.Elements[0].(*object.Atom)
		if !ok {
			return newError(ctx.Line, ctx.Column, "first element of tuple must be an atom, got %s", tuple.Elements[0].Type())
		}

		accumulator = tuple.Elements[1]

		if action.Value == "halt" {
			break
		} else if action.Value != "cont" {
			return newError(ctx.Line, ctx.Column, "first element of tuple must be :cont or :halt, got %s", action.Value)
		}
	}

	return accumulator
}

func iter(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=2", len(args))
	}

	list, ok := args[0].(*object.Array)
	if !ok {
		return newError(ctx.Line, ctx.Column, "first argument to `iter` must be ARRAY, got %s", args[0].Type())
	}

	var fn object.Object
	switch arg := args[1].(type) {
	case *object.Function:
		fn = arg
	case *object.Builtin:
		fn = arg
	default:
		return newError(ctx.Line, ctx.Column, "second argument to `iter` must be FUNCTION, got %s", args[1].Type())
	}

	for i := 0; i < len(list.Elements); i++ {
		var result object.Object
		if function, ok := fn.(*object.Function); ok {
			result = applyFunction(function, []object.Object{list.Elements[i]}, ctx)
		} else if builtin, ok := fn.(*object.Builtin); ok {
			result = builtin.Fn(ctx, []object.Object{list.Elements[i]}...)
		}

		if result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return constants.OK
}

func reduce(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 3 {
		return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=2 or 3", len(args))
	}

	list, ok := args[0].(*object.Array)
	if !ok {
		return newError(ctx.Line, ctx.Column, "first argument to `reduce` must be ARRAY, got %s", args[0].Type())
	}

	var initial object.Object
	var fn *object.Function
	var startIndex int

	if len(args) == 2 {
		if len(list.Elements) == 0 {
			return newError(ctx.Line, ctx.Column, "cannot reduce empty array without initial value")
		}
		initial = list.Elements[0]
		fn, ok = args[1].(*object.Function)
		if !ok {
			return newError(ctx.Line, ctx.Column, "second argument to `reduce` must be FUNCTION, got %s", args[1].Type())
		}
		startIndex = 1
	} else {
		initial = args[1]
		fn, ok = args[2].(*object.Function)
		if !ok {
			return newError(ctx.Line, ctx.Column, "third argument to `reduce` must be FUNCTION, got %s", args[2].Type())
		}
		startIndex = 0
	}

	accumulator := initial
	for i := startIndex; i < len(list.Elements); i++ {
		elem := list.Elements[i]
		result := applyFunction(fn, []object.Object{accumulator, elem}, ctx)
		if result.Type() == object.ERROR_OBJ {
			return result
		}
		accumulator = result
	}

	return accumulator
}

func loop(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want =1", len(args))
	}

	fnObj, ok := args[1].(*object.Function)
	if !ok {
		return newError(ctx.Line, ctx.Column, "second argument to loop must be a function, got %s", args[1].Type())
	}

	accumulator := args[0]

	for {
		result := applyFunction(fnObj, []object.Object{accumulator}, ctx)

		if result.Type() == object.ERROR_OBJ {
			return result
		}

		tuple, ok := result.(*object.Tuple)

		if !ok || len(tuple.Elements) != 2 {
			return newError(ctx.Line, ctx.Column, "function must return a tuple of (:cont, acc) or (:halt, acc)")
		}

		action := tuple.Elements[0].(*object.Atom)
		if !ok {
			return newError(ctx.Line, ctx.Column, "first element of tuple must be an atom, got %s", tuple.Elements[0].Type())
		}

		accumulator = tuple.Elements[1]

		if action.Value == "halt" {
			break
		} else if action.Value != "cont" {
			return newError(ctx.Line, ctx.Column, "first element of tuple must be :cont or :halt, got %s", action.Value)
		}
	}

	return accumulator
}

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
			case *object.Tuple:
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
	"test": {
		Fn: func(ctx *object.EvalContext, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.BOOLEAN_OBJ {
				return newError(ctx.Line, ctx.Column, "first argument to `test` must be BOOLEAN, got %s", args[0].Type())
			}

			if args[1].Type() != object.STRING_OBJ {
				return newError(ctx.Line, ctx.Column, "second argument to `test` must be STRING, got %s", args[1].Type())
			}

			if args[0].(*object.Boolean) == constants.FALSE {
				fmt.Printf("Test failed: %s\n", args[1].(*object.String).Value)
			}
			return constants.OK
		},
	},
}
