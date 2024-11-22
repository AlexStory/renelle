// builtins/map.go

package hostlib

import (
	"renelle/constants"
	"renelle/object"
)

func MapHasKey(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "haskey() takes exactly 2 arguments"}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "haskey() requires a map"}
	}

	key := args[1]

	_, ok = m.Store.Get(key)
	if ok {
		return constants.TRUE
	}

	return constants.FALSE
}

func MapKeys(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "keys() takes exactly 1 argument"}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "keys() requires a map"}
	}

	keys := m.Store.Keys()
	arr := &object.Array{Elements: keys}

	return arr
}

func MapLength(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "length() takes exactly 1 argument"}
	}

	m, ok := args[0].(*object.Map)
	if !ok {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "length() requires a map"}
	}

	return &object.Integer{Value: int64(m.Store.Length)}
}
