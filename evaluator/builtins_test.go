package evaluator

import (
	"renelle/object"
	"testing"
)

func TestBuiltinTypeFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`type("Hello, World!")`, "STRING"},
		{`type(123)`, "INTEGER"},
		{`type(true)`, "BOOLEAN"},
		{`type(3.14)`, "FLOAT"},
		{`type([1, 2, 3])`, "ARRAY"},
		{`type((1, 2, 3))`, "TUPLE"},
		{`type({ a: 1, "b" = 2 })`, "MAP"},
		{`type(:a)`, "ATOM"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		strObj, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if strObj.Value != tt.expected {
			t.Errorf("wrong value. expected=%q, got=%q", tt.expected, strObj.Value)
		}
	}
}
