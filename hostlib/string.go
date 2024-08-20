package hostlib

import (
	"regexp"
	"strings"

	"renelle/constants"
	"renelle/object"
)

// StringConcat concatenates two strings.
func StringConcat(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "concat() takes exactly 2 arguments"}
	}

	str1, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "concat() requires a string"}
	}

	str2, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "concat() requires a string"}
	}

	return &object.String{Value: str1.Value + str2.Value}
}

// StringContains checks if a string contains a substring.
func StringContains(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "contains() takes exactly 2 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "contains() requires a string"}
	}

	substr, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "contains() requires a string"}
	}

	return &object.Boolean{Value: strings.Contains(str.Value, substr.Value)}
}

// StringEndsWith checks if a string ends with a substring.
func StringEndsWith(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "endsWith() takes exactly 2 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "endsWith() requires a string"}
	}

	substr, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "endsWith() requires a string"}
	}

	return &object.Boolean{Value: strings.HasSuffix(str.Value, substr.Value)}
}

// StringIndexOf returns the index of a substring in a string.
func StringIndexOf(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "indexOf() takes exactly 2 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "indexOf() requires a string"}
	}

	substr, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "indexOf() requires a string"}
	}

	return &object.Integer{Value: int64(strings.Index(str.Value, substr.Value))}
}

// StringLength returns the length of a string.
func StringLength(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "length() takes exactly 1 argument"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "length() requires a string"}
	}

	return &object.Integer{Value: int64(len(str.Value))}
}

// StringLower converts a string to lowercase.
func StringLower(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "lower() takes exactly 1 argument"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "lower() requires a string"}
	}

	return &object.String{Value: strings.ToLower(str.Value)}
}

// Returns whether a string matches a regular expression.
func StringMatch(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "match() takes exactly 2 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "match() requires a string"}
	}

	re, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "match() requires a string"}
	}

	matched, err := regexp.MatchString(re.Value, str.Value)
	if err != nil {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: err.Error()}
	}

	if matched {
		return constants.TRUE
	} else {
		return constants.FALSE
	}
}

// StringPadLeft pads a string on the left. using the given character or " " by default
func StringPadLeft(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) < 2 || len(args) > 3 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padLeft() takes 2 or 3 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padLeft() requires a string"}
	}

	length, ok := args[1].(*object.Integer)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padLeft() requires an integer"}
	}

	pad := " "
	if len(args) == 3 {
		padObj, ok := args[2].(*object.String)
		if !ok {
			return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padLeft() requires a string"}
		}
		pad = padObj.Value
	}

	return &object.String{Value: strings.Repeat(pad, int(length.Value)) + str.Value}
}

// StringPadRight pads a string on the right. using the given character or " " by default
func StringPadRight(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) < 2 || len(args) > 3 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padRight() takes 2 or 3 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padRight() requires a string"}
	}

	length, ok := args[1].(*object.Integer)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padRight() requires an integer"}
	}

	pad := " "
	if len(args) == 3 {
		padObj, ok := args[2].(*object.String)
		if !ok {
			return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "padRight() requires a string"}
		}
		pad = padObj.Value
	}

	return &object.String{Value: str.Value + strings.Repeat(pad, int(length.Value))}
}

// StringReplace replaces a substring in a string.
func StringReplace(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 3 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replace() takes exactly 3 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replace() requires a string"}
	}

	old, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replace() requires a string"}
	}

	new, ok := args[2].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replace() requires a string"}
	}

	return &object.String{Value: strings.Replace(str.Value, old.Value, new.Value, 1)}
}

// Replaces all occurrences of a substring in a string.
func StringReplaceAll(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 3 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replaceAll() takes exactly 3 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replaceAll() requires a string"}
	}

	old, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replaceAll() requires a string"}
	}

	new, ok := args[2].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "replaceAll() requires a string"}
	}

	return &object.String{Value: strings.ReplaceAll(str.Value, old.Value, new.Value)}
}

// StringSplit splits a string by a separator.
func StringSplit(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "split() takes 1 or 2 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "split() requires a string"}
	}

	sep := ""
	if len(args) == 2 {
		sepObj, ok := args[1].(*object.String)
		if !ok {
			return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "split() requires a string as a separator"}
		}
		sep = sepObj.Value
	}
	elements := make([]object.Object, 0)
	for _, s := range strings.Split(str.Value, sep) {
		elements = append(elements, &object.String{Value: s})
	}

	return &object.Array{Elements: elements}
}

// StringStartsWith checks if a string starts with a substring.
func StringStartsWith(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "startsWith() takes exactly 2 arguments"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "startsWith() requires a string"}
	}

	substr, ok := args[1].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "startsWith() requires a string"}
	}

	return &object.Boolean{Value: strings.HasPrefix(str.Value, substr.Value)}
}

// StringTrim trims a string.
func StringTrim(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "trim() takes exactly 1 argument"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "trim() requires a string"}
	}

	return &object.String{Value: strings.TrimSpace(str.Value)}
}

// trimEnd trims the end of a string.
func StringTrimEnd(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "trimEnd() takes exactly 1 argument"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "trimEnd() requires a string"}
	}

	return &object.String{Value: strings.TrimRightFunc(str.Value, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})}
}

// trimStart trims the start of a string.
func StringTrimStart(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "trimStart() takes exactly 1 argument"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "trimStart() requires a string"}
	}

	return &object.String{Value: strings.TrimLeftFunc(str.Value, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})}
}

// StringUpper converts a string to uppercase.
func StringUpper(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "upper() takes exactly 1 argument"}
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Line: ctx.Line, Column: ctx.Column, Message: "upper() requires a string"}
	}

	return &object.String{Value: strings.ToUpper(str.Value)}
}
