// object/object.go

package object

import (
	"bytes"
	"fmt"
	"renelle/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	ATOM_OBJ         = "ATOM"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return "\"" + s.Value + "\"" }
func (s *String) Type() ObjectType { return STRING_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Atom struct {
	Value string
}

func (a *Atom) Inspect() string  { return a.Value }
func (a *Atom) Type() ObjectType { return ATOM_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("Line: %d, Column %d: ERROR: %s", e.Line, e.Column, e.Message)
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n\t")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()

}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

type BuiltinFunction func(line, col int, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
