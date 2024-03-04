// builtins/array.go

package hostlib

import (
	"slices"

	"renelle/object"
)

func ArrayReverse(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "reverse() takes exactly 1 argument"}
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "reverse() requires an array"}
	}

	el := arr.Elements
	newArray := make([]object.Object, len(el))
	copy(newArray, el)
	slices.Reverse(newArray)

	return &object.Array{Elements: newArray}
}
