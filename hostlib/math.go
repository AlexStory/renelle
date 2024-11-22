package hostlib

import (
	"math"
	"renelle/object"
)

func MathAbs(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "abs() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		if num.Value < 0 {
			return &object.Integer{Value: -num.Value}
		}
		return num
	case *object.Float:
		if num.Value < 0 {
			return &object.Float{Value: -num.Value}
		}
		return num
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "abs() requires a number"}
	}
}

func MathCeil(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "ceil() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		return num
	case *object.Float:
		if num.Value-float64(int64(num.Value)) > 0 {
			return &object.Integer{Value: int64(num.Value) + 1}
		}
		return &object.Integer{Value: int64(num.Value)}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "ceil() requires a number"}
	}
}

func MathCos(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "cos() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		return &object.Float{Value: math.Cos(float64(num.Value))}
	case *object.Float:
		return &object.Float{Value: math.Cos(num.Value)}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "cos() requires a number"}
	}
}

func MathFloor(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "floor() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		return num
	case *object.Float:
		return &object.Integer{Value: int64(num.Value)}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "floor() requires a number"}
	}
}

func MathMax(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "max() requires exactly 2 arguments"}
	}

	switch num1 := args[0].(type) {
	case *object.Integer:
		switch num2 := args[1].(type) {
		case *object.Integer:
			if num1.Value > num2.Value {
				return num1
			}
			return num2
		case *object.Float:
			if float64(num1.Value) > num2.Value {
				return num1
			}
			return num2
		default:
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "max() requires numbers"}
		}
	case *object.Float:
		switch num2 := args[1].(type) {
		case *object.Integer:
			if num1.Value > float64(num2.Value) {
				return num1
			}
			return num2
		case *object.Float:
			if num1.Value > num2.Value {
				return num1
			}
			return num2
		default:
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "max() requires numbers"}
		}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "max() requires numbers"}
	}
}

func MathMin(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "min() requires exactly 2 arguments"}
	}

	switch num1 := args[0].(type) {
	case *object.Integer:
		switch num2 := args[1].(type) {
		case *object.Integer:
			if num1.Value < num2.Value {
				return num1
			}
			return num2
		case *object.Float:
			if float64(num1.Value) < num2.Value {
				return num1
			}
			return num2
		default:
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "min() requires numbers"}
		}
	case *object.Float:
		switch num2 := args[1].(type) {
		case *object.Integer:
			if num1.Value < float64(num2.Value) {
				return num1
			}
			return num2
		case *object.Float:
			if num1.Value < num2.Value {
				return num1
			}
			return num2
		default:
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "min() requires numbers"}
		}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "min() requires numbers"}
	}
}

func MathPi(ctx *object.EvalContext, args ...object.Object) object.Object {
	return &object.Float{Value: math.Pi}
}

func MathRound(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "round() takes 1 or 2 arguments"}
	}

	precision := 0
	if len(args) == 2 {
		if prec, ok := args[1].(*object.Integer); ok {
			precision = int(prec.Value)
		} else {
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "precision must be an integer"}
		}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		return num
	case *object.Float:
		factor := math.Pow(10, float64(precision))
		roundedValue := math.Round(num.Value*factor) / factor
		if precision == 0 {
			return &object.Integer{Value: int64(roundedValue)}
		}
		return &object.Float{Value: roundedValue}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "round() requires a number"}
	}
}

func MathSin(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "sin() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		return &object.Float{Value: math.Sin(float64(num.Value))}
	case *object.Float:
		return &object.Float{Value: math.Sin(num.Value)}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "sin() requires a number"}
	}
}

func MathSqrt(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "sqrt() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		if num.Value < 0 {
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "sqrt() requires a non-negative number"}
		}
		return &object.Float{Value: math.Sqrt(float64(num.Value))}
	case *object.Float:
		if num.Value < 0 {
			return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "sqrt() requires a non-negative number"}
		}
		return &object.Float{Value: math.Sqrt(num.Value)}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "sqrt() requires a number"}
	}
}

func MathTan(ctx *object.EvalContext, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "tan() takes exactly 1 argument"}
	}

	switch num := args[0].(type) {
	case *object.Integer:
		return &object.Float{Value: math.Tan(float64(num.Value))}
	case *object.Float:
		return &object.Float{Value: math.Tan(num.Value)}
	default:
		return &object.Error{FileName: ctx.FileName, Line: ctx.Line, Column: ctx.Column, Message: "tan() requires a number"}
	}
}
