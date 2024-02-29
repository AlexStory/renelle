// evaluator/evaluator.go

package evaluator

import (
	"fmt"
	"math"

	"renelle/ast"
	"renelle/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL   = &object.Atom{Value: "nil"}
	OK    = &object.Atom{Value: "ok"}
)

var atoms = map[string]*object.Atom{
	"nil": NIL,
	"ok":  OK,
}

func Eval(node ast.Node, env *object.Environment, ctx *object.EvalContext) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		result := evalProgram(node.Statements, env, ctx)

		mainFunc, ok := env.Get("main")
		if ok {
			return applyFunction(mainFunc, []object.Object{}, ctx)
		}

		return result

	case *ast.FunctionStatement:
		env.Set(node.Name.Value, &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env})

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env, ctx)

	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env, ctx)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env, ctx)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env, ctx)
		if isError(val) {
			return val
		}
		switch left := node.Left.(type) {
		case *ast.Identifier:
			if left.Value != "_" {
				env.Set(left.Value, val)
			}
		case *ast.TupleLiteral:
			ctx.Line = node.Token.Line
			ctx.Column = node.Token.Column
			return handleTupleDestructuring(left, val, env, ctx)
		case *ast.ArrayLiteral:
			ctx.Line = node.Token.Line
			ctx.Column = node.Token.Column
			return handleArrayDestructuring(left, val, env, ctx)
		default:
			return newError(node.Token.Line, node.Token.Column, "invalid left-hand side of assignment")
		}
		return OK

	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, ctx)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.TupleLiteral:
		elements := evalExpressions(node.Elements, env, ctx)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Tuple{Elements: elements}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.AtomLiteral:
		return getOrCreateAtom(node.Value)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env, ctx)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Line, node.Token.Column)

	case *ast.InfixExpression:
		if node.Operator == "|>" {
			if right, ok := node.Right.(*ast.CallExpression); !ok {
				return newError(node.Token.Line, node.Token.Column, "pipe operator must be followed by a function call")
			} else {
				right.Arguments = append([]ast.Expression{node.Left}, right.Arguments...)
				return Eval(right, env, ctx)
			}
		}

		left := Eval(node.Left, env, ctx)
		if isError(left) {
			return left
		}

		if node.Operator == "and" && !isTruthy(left) {
			return left
		}

		if node.Operator == "or" && isTruthy(left) {
			return left
		}

		right := Eval(node.Right, env, ctx)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token.Line, node.Token.Column)

	case *ast.IfExpression:
		return evalIfExpression(node, env, ctx)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.IndexExpression:
		left := Eval(node.Left, env, ctx)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env, ctx)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index, node.Token.Line, node.Token.Column)

	case *ast.CallExpression:
		function := Eval(node.Function, env, ctx)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env, ctx)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, ctx)
	}

	return nil
}

func evalExpressions(exps []ast.Expression, env *object.Environment, ctx *object.EvalContext) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env, ctx)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func handleTupleDestructuring(tuple *ast.TupleLiteral, val object.Object, env *object.Environment, ctx *object.EvalContext) object.Object {
	tupleObject, ok := val.(*object.Tuple)
	if !ok {
		return newError(ctx.Line, ctx.Column, "right-hand side of assignment is not a tuple")
	}
	if len(tuple.Elements) != len(tupleObject.Elements) {
		return newError(ctx.Line, ctx.Column, "cannot destructure tuple: size mismatch")
	}
	for i, el := range tuple.Elements {
		switch el := el.(type) {
		case *ast.Identifier:
			if el.Value == "_" {
				continue
			}
			env.Set(el.Value, tupleObject.Elements[i])
		default:
			leftVal := Eval(el, env, ctx)
			if isError(leftVal) {
				return leftVal
			}
			if !objectEquals(leftVal, tupleObject.Elements[i]) {
				tok := tuple.Elements[i]
				return newError(tok.T().Line, tok.T().Column, "cannot destructure tuple: value mismatch")
			}
		}
	}
	return OK
}

func handleArrayDestructuring(array *ast.ArrayLiteral, val object.Object, env *object.Environment, ctx *object.EvalContext) object.Object {
	arrayObject, ok := val.(*object.Array)
	if !ok {
		return newError(ctx.Line, ctx.Column, "right-hand side of assignment is not an array")
	}
	if len(array.Elements) != len(arrayObject.Elements) {
		return newError(ctx.Line, ctx.Column, "cannot destructure array: size mismatch")
	}
	for i, el := range array.Elements {
		switch el := el.(type) {
		case *ast.Identifier:
			if el.Value == "_" {
				continue
			}
			env.Set(el.Value, arrayObject.Elements[i])
		default:
			leftVal := Eval(el, env, ctx)
			if isError(leftVal) {
				return leftVal
			}
			if !objectEquals(leftVal, arrayObject.Elements[i]) {
				tok := array.Elements[i]
				return newError(tok.T().Line, tok.T().Column, "cannot destructure array: value mismatch")
			}
		}
	}
	return OK
}

func applyFunction(fn object.Object, args []object.Object, ctx *object.EvalContext) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args, ctx)
		evaluated := Eval(fn.Body, extendedEnv, ctx)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(ctx, args...)
	default:
		return newError(ctx.Line, ctx.Column, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object, ctx *object.EvalContext) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalProgram(stmts []ast.Statement, env *object.Environment, ctx *object.EvalContext) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env, ctx)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment, ctx *object.EvalContext) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env, ctx)

		rt := result.Type()
		if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
			return result
		}
	}

	return result

}

func evalIndexExpression(left, index object.Object, line, col int) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, line, col)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, line, col)
	case left.Type() == object.TUPLE_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalTupleIndexExpression(left, index, line, col)
	default:
		return newError(line, col, "index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, line, col int) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NIL
	}

	return arrayObject.Elements[idx]
}

func evalTupleIndexExpression(tuple, index object.Object, line, col int) object.Object {
	tupleObject := tuple.(*object.Tuple)
	idx := index.(*object.Integer).Value
	max := int64(len(tupleObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NIL
	}

	return tupleObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, line, col int) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return NIL
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError(node.Token.Line, node.Token.Column, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(operator string, right object.Object, line, col int) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, line, col)
	default:
		return newError(line, col, "unknown operator: %s%s", operator, right.Type())

	}
}

func evalInfixExpression(operator string, left, right object.Object, line, col int) object.Object {
	switch {
	case operator == "and":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "or":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, line, col)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatInfixExpression(operator, left, &object.Float{Value: float64(right.(*object.Integer).Value)})
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, &object.Float{Value: float64(left.(*object.Integer).Value)}, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, line, col)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError(line, col, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(line, col, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, line, col int) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "**":
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(line, col, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: float64(int64(leftVal) % int64(rightVal))}
	case "**":
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NIL
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, line, col int) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(line, col, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment, ctx *object.EvalContext) object.Object {
	condition := Eval(ie.Condition, env, ctx)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env, ctx)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env, ctx)
	} else {
		return NIL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, line, col int) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	default:
		return newError(line, col, "unknown operator: -%s", right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func getOrCreateAtom(value string) *object.Atom {
	if atom, ok := atoms[value]; ok {
		return atom
	}
	atom := &object.Atom{Value: value}
	atoms[value] = atom
	return atom
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(line, column int, format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...), Line: line, Column: column}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func objectEquals(a, b object.Object) bool {
	switch a := a.(type) {
	case *object.Integer:
		b, ok := b.(*object.Integer)
		return ok && a.Value == b.Value
	case *object.Float:
		b, ok := b.(*object.Float)
		return ok && a.Value == b.Value
	case *object.String:
		b, ok := b.(*object.String)
		return ok && a.Value == b.Value
	case *object.Boolean:
		b, ok := b.(*object.Boolean)
		return ok && a.Value == b.Value
	case *object.Atom:
		b, ok := b.(*object.Atom)
		return ok && a.Value == b.Value
	case *object.ReturnValue:
		b, ok := b.(*object.ReturnValue)
		return ok && objectEquals(a.Value, b.Value)
	case *object.Array:
		b, ok := b.(*object.Array)
		if !ok || len(a.Elements) != len(b.Elements) {
			return false
		}
		for i, el := range a.Elements {
			if !objectEquals(el, b.Elements[i]) {
				return false
			}
		}
		return true
	case *object.Tuple:
		b, ok := b.(*object.Tuple)
		if !ok || len(a.Elements) != len(b.Elements) {
			return false
		}
		for i, el := range a.Elements {
			if !objectEquals(el, b.Elements[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
