package evaluator

import "renelle/object"

func VectorVectorAdd(ctx *object.EvalContext, left *object.Array, right *object.Array) object.Object {
	if len(left.Elements) != len(right.Elements) {
		return newError(ctx.Line, ctx.Column, "vector length mismatch: %d != %d", len(left.Elements), len(right.Elements))
	}

	elements := make([]object.Object, len(left.Elements))
	for i := range left.Elements {
		elements[i] = evalInfixExpression(ctx, "+", left.Elements[i], right.Elements[i], ctx.Line, ctx.Column)
	}

	return &object.Array{Elements: elements}
}

func VectorScalarAdd(ctx *object.EvalContext, left *object.Array, right object.Object) object.Object {
	elements := make([]object.Object, len(left.Elements))
	for i := range left.Elements {
		elements[i] = evalInfixExpression(ctx, "+", left.Elements[i], right, ctx.Line, ctx.Column)
	}

	return &object.Array{Elements: elements}
}
