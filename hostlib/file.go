// hostlib/file.go

package hostlib

import (
	"os"

	"renelle/constants"
	"renelle/object"
)

func FileOpen(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "open() takes exactly 1 argument"}
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "open() requires a string"}
	}

	file, err := os.ReadFile(path.Value)

	res := &object.Tuple{}
	elements := make([]object.Object, 0)

	if err != nil {
		elements = append(elements, constants.ERROR)
		elements = append(elements, &object.String{Value: err.Error()})
	} else {
		elements = append(elements, constants.OK)
		elements = append(elements, &object.String{Value: string(file)})
	}

	res.Elements = elements
	return res
}

func FileOpenBang(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "open!() takes exactly 1 argument"}
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "open!() requires a string"}
	}

	file, err := os.ReadFile(path.Value)
	if err != nil {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: err.Error()}
	}

	return &object.String{Value: string(file)}
}

func FileWrite(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "write() takes exactly 2 arguments"}
	}

	content, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "write() requires a string"}
	}

	path, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "write() requires a string"}
	}

	res := &object.Tuple{}
	elements := make([]object.Object, 0)

	err := os.WriteFile(path.Value, []byte(content.Value), 0644)
	if err != nil {
		elements = append(elements, constants.ERROR)
		elements = append(elements, &object.String{Value: err.Error()})
	} else {
		elements = append(elements, constants.OK)
		elements = append(elements, constants.NIL)
	}

	res.Elements = elements
	return res
}

func FileWriteBang(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "write!() takes exactly 2 arguments"}
	}

	content, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "write!() requires a string"}
	}

	path, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "write!() requires a string"}
	}

	err := os.WriteFile(path.Value, []byte(content.Value), 0644)
	if err != nil {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: err.Error()}
	}

	return constants.OK
}
