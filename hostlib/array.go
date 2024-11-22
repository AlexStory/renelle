// builtins/array.go

package hostlib

import (
	"slices"

	"renelle/object"
)

func ArrayReverse(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "reverse() takes exactly 1 argument"}
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "reverse() requires an array"}
	}

	el := arr.Elements
	newArray := make([]object.Object, len(el))
	copy(newArray, el)
	slices.Reverse(newArray)

	return &object.Array{Elements: newArray}
}

func ArrayRange(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "range() takes 1 or 2 arguments"}
	}

	var start, stop int64

	if len(args) == 1 {
		stopArg, ok := args[0].(*object.Integer)
		if !ok {
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "range() requires an integer"}
		}
		start = 0
		stop = stopArg.Value
	} else {
		startArg, ok := args[0].(*object.Integer)
		if !ok {
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "range() requires integer arguments"}
		}
		stopArg, ok := args[1].(*object.Integer)
		if !ok {
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "range() requires integer arguments"}
		}
		start = startArg.Value
		stop = stopArg.Value
	}

	if start >= stop {
		return &object.Array{Elements: []object.Object{}}
	}

	elements := make([]object.Object, stop-start)
	for i := start; i < stop; i++ {
		elements[i-start] = &object.Integer{Value: i}
	}

	return &object.Array{Elements: elements}

}
