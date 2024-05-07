// builtins/map.go

package hostlib

import (
	"renelle/object"
)

func MapLength(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "length() takes exactly 1 argument"}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "length() requires a map"}
	}

	return &object.Integer{Value: int64(m.Store.Length)}
}
