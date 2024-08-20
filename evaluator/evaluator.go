// evaluator/evaluator.go

package evaluator

import (
	"embed"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"renelle/ast"
	"renelle/constants"
	"renelle/hostlib"
	"renelle/lexer"
	"renelle/object"
	"renelle/parser"
	"renelle/stdlib"
)

var atoms = map[string]*object.Atom{
	"nil":   constants.NIL,
	"ok":    constants.OK,
	"error": constants.ERROR,
}

func ApplyFunction(fn object.Object, args []object.Object, ctx *object.EvalContext) object.Object {
	return applyFunction(fn, args, ctx)
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
				if unicode.IsUpper(rune(left.Value[0])) {
					return newError(node.Token.Line, node.Token.Column, "local variables can not start with an uppercase letter")
				}
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
		case *ast.MapLiteral:
			ctx.Line = node.Token.Line
			ctx.Column = node.Token.Column
			return handleMapDestructuring(left, val, env, ctx)
		default:
			return newError(node.Token.Line, node.Token.Column, "invalid left-hand side of assignment")
		}
		return val
	case *ast.Module:
		moduleEnv := object.NewEnclosedEnvironment(env)
		for _, statement := range node.Body {
			ret := Eval(statement, moduleEnv, ctx)
			if isError(ret) {
				return ret
			}
		}
		module := &object.Module{Name: node.Name.Value, Environment: moduleEnv}
		env.SetModule(node.Name.Value, module)
		return module

	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.InterpolatedStringLiteral:
		var sb strings.Builder
		for _, part := range node.Segments {
			evaluated := Eval(part, env, ctx)
			if isError(evaluated) {
				return evaluated
			}
			switch evaluated := evaluated.(type) {
			case *object.String:
				sb.WriteString(evaluated.Value)
			case *object.Integer:
				sb.WriteString(evaluated.Inspect())
			case *object.Float:
				sb.WriteString(evaluated.Inspect())
			default:
				sb.WriteString(evaluated.Inspect())
			}
		}
		return &object.String{Value: sb.String()}

	case *ast.MapLiteral:
		return evalMapLiteral(node, env, ctx)
	case *ast.MapUpdateLiteral:
		return evalMapUpdateLiteral(node, env, ctx)
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
	case *ast.PropertyAccessExpression:
		left := Eval(node.Left, env, ctx)
		if isError(left) {
			return left
		}
		ctx.Line = node.Token.Line
		ctx.Column = node.Token.Column
		return evalPropertyAccessExpression(left, node.Right, env, ctx)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env, ctx)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Line, node.Token.Column)

	case *ast.InfixExpression:
		if node.Operator == "|>" {
			switch right := node.Right.(type) {
			case *ast.CallExpression:
				right.Arguments = append([]ast.Expression{node.Left}, right.Arguments...)
				return Eval(right, env, ctx)
			case *ast.PropertyAccessExpression:
				// Assuming that the property access expression has a CallExpression as its property
				if callExpr, ok := right.Right.(*ast.CallExpression); ok {
					callExpr.Arguments = append([]ast.Expression{node.Left}, callExpr.Arguments...)
					return Eval(right, env, ctx)
				} else {
					return newError(node.Token.Line, node.Token.Column, "pipe operator must be followed by a function call")
				}
			case *ast.FunctionLiteral:
				if len(right.Parameters) != 1 {
					return newError(node.Token.Line, node.Token.Column, "function literal must take exactly one argument")
				}
				newCall := &ast.CallExpression{
					Function:  right,
					Arguments: []ast.Expression{node.Left},
				}
				return Eval(newCall, env, ctx)
			default:
				return newError(node.Token.Line, node.Token.Column, "pipe operator must be followed by a function call")
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

	case *ast.CaseExpression:
		testVal := Eval(node.Test, env, ctx)
		if isError(testVal) {
			return testVal
		}
		for i, condition := range node.Conditions {
			switch condition := condition.(type) {
			case *ast.Identifier:
				if condition.Value == "_" {
					newEnv := object.NewEnclosedEnvironment(env)
					return Eval(node.Consequences[i], newEnv, ctx)
				}
			case *ast.TupleLiteral:
				ctx.Line = node.Token.Line
				ctx.Column = node.Token.Column
				newEnv := object.NewEnclosedEnvironment(env)
				err := handleTupleDestructuring(condition, testVal, newEnv, ctx)
				if isError(err) {
					continue
				}
				return Eval(node.Consequences[i], newEnv, ctx)
			case *ast.ArrayLiteral:
				ctx.Line = node.Token.Line
				ctx.Column = node.Token.Column
				newEnv := object.NewEnclosedEnvironment(env)
				err := handleArrayDestructuring(condition, testVal, newEnv, ctx)
				if isError(err) {
					continue
				}
				return Eval(node.Consequences[i], newEnv, ctx)
			default:
				conditionVal := Eval(condition, env, ctx)
				if isError(conditionVal) {
					return conditionVal
				}
				if object.Equals(conditionVal, testVal) {
					newEnv := object.NewEnclosedEnvironment(env)
					return Eval(node.Consequences[i], newEnv, ctx)
				}
			}
		}
		return newError(node.Token.Line, node.Token.Column, "no matching case")
	case *ast.IfExpression:
		return evalIfExpression(node, env, ctx)

	case *ast.CondExpression:
		return evalCondExpression(node, env, ctx)

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
			} else if unicode.IsUpper(rune(el.Value[0])) {
				return newError(el.Token.Line, el.Token.Column, "local variables can not start with an uppercase letter")
			}
			env.Set(el.Value, tupleObject.Elements[i])
		case *ast.TupleLiteral:
			return handleTupleDestructuring(el, tupleObject.Elements[i], env, ctx)
		case *ast.ArrayLiteral:
			return handleArrayDestructuring(el, tupleObject.Elements[i], env, ctx)
		default:
			leftVal := Eval(el, env, ctx)
			if isError(leftVal) {
				return leftVal
			}
			if !object.Equals(leftVal, tupleObject.Elements[i]) {
				tok := tuple.Elements[i]
				return newError(tok.T().Line, tok.T().Column, "cannot destructure tuple: value mismatch")
			}
		}
	}
	return constants.OK
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
			} else if unicode.IsUpper(rune(el.Value[0])) {
				return newError(el.Token.Line, el.Token.Column, "local variables can not start with an uppercase letter")
			}
			env.Set(el.Value, arrayObject.Elements[i])
		case *ast.TupleLiteral:
			return handleTupleDestructuring(el, arrayObject.Elements[i], env, ctx)
		case *ast.ArrayLiteral:
			return handleArrayDestructuring(el, arrayObject.Elements[i], env, ctx)
		default:
			leftVal := Eval(el, env, ctx)
			if isError(leftVal) {
				return leftVal
			}
			if !object.Equals(leftVal, arrayObject.Elements[i]) {
				tok := array.Elements[i]
				return newError(tok.T().Line, tok.T().Column, "cannot destructure array: value mismatch")
			}
		}
	}
	return constants.OK
}

func handleMapDestructuring(left *ast.MapLiteral, val object.Object, env *object.Environment, ctx *object.EvalContext) object.Object {
	mapObj, ok := val.(*object.Map)
	if !ok {
		return newError(ctx.Line, ctx.Column, "expected map, got %s", val.Type())
	}

	for keyExpr, valueExpr := range left.Pairs {
		// Evaluate the key expression
		keyVal := Eval(keyExpr, env, ctx)
		if isError(keyVal) {
			return keyVal
		}

		// Get the value from the map
		value, ok := mapObj.Get(keyVal)
		if !ok {
			return newError(ctx.Line, ctx.Column, "key not found: %s", keyVal.Inspect())
		}

		// Handle the value expression based on its type
		switch valueExpr := valueExpr.(type) {
		case *ast.Identifier:
			// Set the value in the environment under the name given by the value expression
			if unicode.IsUpper(rune(valueExpr.Value[0])) {
				return newError(valueExpr.Token.Line, valueExpr.Token.Column, "local variables can not start with an uppercase letter")
			}
			if valueExpr.Value != "_" {
				env.Set(valueExpr.Value, value)
			}
		case *ast.MapLiteral:
			// Recursively handle nested map destructuring
			return handleMapDestructuring(valueExpr, value, env, ctx)
		case *ast.ArrayLiteral:
			// Handle array destructuring
			return handleArrayDestructuring(valueExpr, value, env, ctx)
		case *ast.TupleLiteral:
			// Handle tuple destructuring
			return handleTupleDestructuring(valueExpr, value, env, ctx)
		default:
			// Evaluate the value expression
			valueVal := Eval(valueExpr, env, ctx)
			if isError(valueVal) {
				return valueVal
			}

			// Check if the value matches the value expression
			if !object.Equals(valueVal, value) {
				return newError(ctx.Line, ctx.Column, "cannot destructure map: value mismatch")
			}
		}
	}

	return val
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment, ctx *object.EvalContext) object.Object {
	hashTableSize := int(float64(len(node.Pairs)) / 0.7)
	hashTable := object.NewHashTable(hashTableSize)
	mapObject := &object.Map{Store: hashTable}

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env, ctx)
		if isError(key) {
			return key
		}

		_, ok := key.(object.Hashable)
		if !ok {
			return newError(ctx.Line, ctx.Column, "unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env, ctx)
		if isError(value) {
			return value
		}

		mapObject.Put(key, value)
	}

	return mapObject
}

func evalPropertyAccessExpression(left object.Object, right ast.Expression, env *object.Environment, ctx *object.EvalContext) object.Object {
	switch left := left.(type) {
	case *object.Module:
		switch right := right.(type) {
		case *ast.Identifier:
			// Get the function from the module
			funcObj, ok := left.Environment.Get(right.Value)
			if !ok {
				return newError(ctx.Line, ctx.Column, "property %s not found", right.Value)
			}

			// Check if the object is a function or a built-in
			switch funcObj := funcObj.(type) {
			case *object.Function, *object.Builtin:
				return funcObj
			default:
				return newError(ctx.Line, ctx.Column, "property %s is not a function", right.Value)
			}
		case *ast.CallExpression:
			// Ensure the function expression is an Identifier
			ident, ok := right.Function.(*ast.Identifier)
			if !ok {
				return newError(ctx.Line, ctx.Column, "invalid function call: %s", right.Function.String())
			}

			// Get the function from the module
			funcObj, ok := left.Environment.Get(ident.Value)
			if !ok {
				return newError(ctx.Line, ctx.Column, "function %s not found", ident.Value)
			}

			switch funcObj := funcObj.(type) {
			case *object.Function:
				// Evaluate the arguments
				args := evalExpressions(right.Arguments, env, ctx)

				// Apply the function to the arguments
				return applyFunction(funcObj, args, ctx)
			case *object.Builtin:
				// Evaluate the arguments
				args := evalExpressions(right.Arguments, env, ctx)

				// Apply the built-in function to the arguments
				return funcObj.Fn(ctx, args...)
			default:
				return newError(ctx.Line, ctx.Column, "property %s is not a function", ident.Value)
			}

		default:
			return newError(ctx.Line, ctx.Column, "invalid property access: %s", right.String())
		}
	case *object.Map:
		// Ensure the right expression is an Identifier
		ident, ok := right.(*ast.Identifier)
		if !ok {
			return newError(ctx.Line, ctx.Column, "invalid property access: %s", right.String())
		}

		// Get or create the atom for the identifier
		key := getOrCreateAtom(ident.Value)

		value, ok := left.Get(key)
		if !ok {
			return constants.NIL
		}
		return value
	default:
		return newError(ctx.Line, ctx.Column, "property access not supported: %s", left.Type())
	}
}

func evalMapUpdateLiteral(node *ast.MapUpdateLiteral, env *object.Environment, ctx *object.EvalContext) object.Object {
	mapObj := Eval(node.Left, env, ctx)
	if isError(mapObj) {
		return mapObj
	}

	mapObjTyped, ok := mapObj.(*object.Map)
	if !ok {
		return newError(ctx.Line, ctx.Column, "not a map: %s", mapObj.Type())
	}

	mapCopy := mapObjTyped.Copy(len(node.Right))
	for keyNode, valueNode := range node.Right {
		key := Eval(keyNode, env, ctx)
		if isError(key) {
			return key
		}

		value := Eval(valueNode, env, ctx)
		if isError(value) {
			return value
		}

		mapCopy.Put(key, value)
	}

	return mapCopy
}

func applyFunction(fn object.Object, args []object.Object, ctx *object.EvalContext) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError(ctx.Line, ctx.Column, "wrong number of arguments. got=%d, want=%d", len(args), len(fn.Parameters))
		}
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
	case left.Type() == object.MAP_OBJ:
		return evalMapIndexExpression(left, index, line, col)
	default:
		return newError(line, col, "index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, line, col int) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return constants.NIL
	}

	return arrayObject.Elements[idx]
}

func evalMapIndexExpression(mapObject, index object.Object, line, col int) object.Object {
	mapObj := mapObject.(*object.Map)
	value, ok := mapObj.Get(index)
	if !ok {
		return constants.NIL
	}
	return value
}

func evalTupleIndexExpression(tuple, index object.Object, line, col int) object.Object {
	tupleObject := tuple.(*object.Tuple)
	idx := index.(*object.Integer).Value
	max := int64(len(tupleObject.Elements) - 1)

	if idx < 0 || idx > max {
		return constants.NIL
	}

	return tupleObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, line, col int) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return constants.NIL
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if unicode.IsUpper([]rune(node.Value)[0]) {
		if module, ok := env.GetModule(node.Value); ok {
			return module
		} else {
			module := loadModule(node.Value, env)
			if error, ok := module.(*object.Error); ok {
				return error
			}
			return module
		}
	} else {
		if val, ok := env.Get(node.Value); ok {
			return val
		}
		if val, ok := builtins[node.Value]; ok {
			return val
		}
	}
	return newError(node.Token.Line, node.Token.Column, "identifier not found: "+node.Value)
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
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(operator, left, right, line, col)
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
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
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
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return constants.NIL
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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError(line, col, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalArrayInfixExpression(operator string, left, right object.Object, line, col int) object.Object {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)
	switch operator {
	case "+":
		return &object.Array{Elements: append(leftVal.Elements, rightVal.Elements...)}
	case "==":
		if len(leftVal.Elements) != len(rightVal.Elements) {
			return constants.FALSE
		}
		for i, el := range leftVal.Elements {
			if !object.Equals(el, rightVal.Elements[i]) {
				return constants.FALSE
			}
		}
		return constants.TRUE
	case "!=":
		if len(leftVal.Elements) != len(rightVal.Elements) {
			return constants.TRUE
		}
		for i, el := range leftVal.Elements {
			if !object.Equals(el, rightVal.Elements[i]) {
				return constants.TRUE
			}
		}
		return constants.FALSE
	default:
		return newError(line, col, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case constants.TRUE:
		return constants.FALSE
	case constants.FALSE:
		return constants.TRUE
	case constants.NIL:
		return constants.TRUE
	default:
		return constants.FALSE
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
		return constants.NIL
	}
}

func evalCondExpression(ce *ast.CondExpression, env *object.Environment, ctx *object.EvalContext) object.Object {
	for i, cond := range ce.Conditions {
		condVal := Eval(cond, env, ctx)
		if isError(condVal) {
			return condVal
		}

		if isTruthy(condVal) {
			return evalBlockStatements(ce.Consequences[i].Statements, env, ctx)
		}
	}
	return constants.NIL
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
		return constants.TRUE
	}
	return constants.FALSE
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
	case constants.NIL:
		return false
	case constants.FALSE:
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

func loadModule(moduleName string, env *object.Environment) object.Object {
	moduleParts := strings.Split(moduleName, ".")
	// Convert each part of the module name to lowercase
	for i, part := range moduleParts {
		moduleParts[i] = strings.ToLower(part)
	}
	// Join the parts with the OS-specific path separator and append the file extension
	localModulePath := filepath.Join(append([]string{"src"}, moduleParts[1:]...)...) + ".rnl"
	depsModulePath := filepath.Join(append([]string{".deps", moduleParts[0], "src"}, moduleParts[1:]...)...) + ".rnl"
	stdLibModulePath := filepath.Join(moduleParts...) + ".rnl"
	if _, err := stdlib.Files.Open(stdLibModulePath); err == nil {
		return loadModuleFromEmbedFS(stdlib.Files, stdLibModulePath, env, moduleName)
	}
	// Check the local src directory first
	if _, err := os.Stat(localModulePath); err == nil {
		return loadModuleFromFile(localModulePath, env, moduleName)
	}
	// If the module is not found locally, check the .deps directory
	if _, err := os.Stat(depsModulePath); err == nil {
		return loadModuleFromFile(depsModulePath, env, moduleName)
	}

	return newError(0, 0, "module not found: %s", moduleName)
}

func loadModuleFromFile(modulePath string, env *object.Environment, moduleName string) object.Object {
	moduleContent, err := os.ReadFile(modulePath)
	if err != nil {
		return newError(0, 0, "error reading module file: %s", err)
	}

	l := lexer.New(string(moduleContent))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return newError(0, 0, "parser errors: %v", p.Errors())
	}

	ctx := object.NewEvalContext()
	Eval(program, env.Root(), ctx)
	if module, ok := env.GetModule(moduleName); ok {
		return module
	}

	return newError(0, 0, "unknown error loading module: %s", moduleName)
}

func loadModuleFromEmbedFS(fs embed.FS, modulePath string, env *object.Environment, moduleName string) object.Object {
	moduleContent, err := fs.ReadFile(modulePath)
	if err != nil {
		return newError(0, 0, "error reading module file: %s", err)
	}

	l := lexer.New(string(moduleContent))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return newError(0, 0, "parser errors: %v", p.Errors())
	}

	ctx := object.NewEvalContext()
	Eval(program, env.Root(), ctx)
	if module, ok := env.GetModule(moduleName); ok {
		switch moduleName {
		case "Array":
			module.Environment.Set("iter", &object.Builtin{Fn: iter})
			module.Environment.Set("range", &object.Builtin{Fn: hostlib.ArrayRange})
			module.Environment.Set("reverse", &object.Builtin{Fn: hostlib.ArrayReverse})
			module.Environment.Set("reduce", &object.Builtin{Fn: reduce})
			module.Environment.Set("reduce_while", &object.Builtin{Fn: reduceWhile})
		case "File":
			module.Environment.Set("open", &object.Builtin{Fn: hostlib.FileOpen})
			module.Environment.Set("open!", &object.Builtin{Fn: hostlib.FileOpenBang})
			module.Environment.Set("write", &object.Builtin{Fn: hostlib.FileWrite})
			module.Environment.Set("write!", &object.Builtin{Fn: hostlib.FileWriteBang})
		case "Map":
			module.Environment.Set("has_key?", &object.Builtin{Fn: hostlib.MapHasKey})
			module.Environment.Set("keys", &object.Builtin{Fn: hostlib.MapKeys})
			module.Environment.Set("length", &object.Builtin{Fn: hostlib.MapLength})
		case "Math":
			module.Environment.Set("abs", &object.Builtin{Fn: hostlib.MathAbs})
			module.Environment.Set("ceiling", &object.Builtin{Fn: hostlib.MathCeil})
			module.Environment.Set("cos", &object.Builtin{Fn: hostlib.MathCos})
			module.Environment.Set("floor", &object.Builtin{Fn: hostlib.MathFloor})
			module.Environment.Set("max", &object.Builtin{Fn: hostlib.MathMax})
			module.Environment.Set("min", &object.Builtin{Fn: hostlib.MathMin})
			module.Environment.Set("pi", &object.Builtin{Fn: hostlib.MathPi})
			module.Environment.Set("round", &object.Builtin{Fn: hostlib.MathRound})
			module.Environment.Set("sin", &object.Builtin{Fn: hostlib.MathSin})
			module.Environment.Set("sqrt", &object.Builtin{Fn: hostlib.MathSqrt})
			module.Environment.Set("tan", &object.Builtin{Fn: hostlib.MathTan})
		case "String":
			module.Environment.Set("concat", &object.Builtin{Fn: hostlib.StringConcat})
			module.Environment.Set("contains?", &object.Builtin{Fn: hostlib.StringContains})
			module.Environment.Set("ends_with?", &object.Builtin{Fn: hostlib.StringEndsWith})
			module.Environment.Set("index_of", &object.Builtin{Fn: hostlib.StringIndexOf})
			module.Environment.Set("length", &object.Builtin{Fn: hostlib.StringLength})
			module.Environment.Set("lower", &object.Builtin{Fn: hostlib.StringLower})
			module.Environment.Set("match?", &object.Builtin{Fn: hostlib.StringMatch})
			module.Environment.Set("pad_left", &object.Builtin{Fn: hostlib.StringPadLeft})
			module.Environment.Set("pad_right", &object.Builtin{Fn: hostlib.StringPadRight})
			module.Environment.Set("replace", &object.Builtin{Fn: hostlib.StringReplace})
			module.Environment.Set("replace_all", &object.Builtin{Fn: hostlib.StringReplaceAll})
			module.Environment.Set("split", &object.Builtin{Fn: hostlib.StringSplit})
			module.Environment.Set("starts_with?", &object.Builtin{Fn: hostlib.StringStartsWith})
			module.Environment.Set("trim", &object.Builtin{Fn: hostlib.StringTrim})
			module.Environment.Set("trim_end", &object.Builtin{Fn: hostlib.StringTrimEnd})
			module.Environment.Set("trim_start", &object.Builtin{Fn: hostlib.StringTrimStart})
			module.Environment.Set("upper", &object.Builtin{Fn: hostlib.StringUpper})
		}
		return module
	}

	return newError(0, 0, "unknown error loading module: %s", moduleName)
}
