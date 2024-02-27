// evaluator/evaluator_test.go

package evaluator

import (
	"fmt"
	"renelle/lexer"
	"renelle/object"
	"renelle/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"10.0", 10.0},
		{"-3.14", -3.14},
		{"-10.0", -10.0},
		{"2.0 * 2.0 * 2.0 * 2.0 * 2.0", 32.0},
		{"-50.0 + 100.0 + -50.0", 0.0},
		{"5.0 * 2.0 + 10.0", 20.0},
		{"5.0 + 2.0 * 10.0", 25.0},
		{"20.0 + 2.0 * -10.0", 0.0},
		{"50.0 / 2.0 * 2.0 + 10.0", 60.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestMixedMath(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14 + 5", 8.14},
		{"3.0 - 5", -2.0},
		{"3.0 * 5", 15.0},
		{"3.14 / 5", 0.628},
		{"3.14 % 5", 3.00},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)

	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"3 < 5", true},
		{"5 < 3", false},
		{"3 > 5", false},
		{"5 > 3", true},
		{"3 == 3", true},
		{"3 != 3", false},
		{"3 == 5", false},
		{"3 != 5", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"true and true", true},
		{"true and false", false},
		{"true or false", true},
		{"(1 < 2) == true", true},
		{"false or false and true", false},
		{`"hello" == "hello"`, true},
		{`"hello" == "world"`, false},
		{`"hello" != "world"`, true},
		{`"hello" != "hello"`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalAtomExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{":foo", ":foo"},
		{":bar", ":bar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testAtomObject(t, evaluated, tt.expected)
	}
}

func testAtomObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Atom)
	if !ok {
		t.Errorf("object is not Atom. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!:nil", true},
		{"!:ok", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != NIL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10", 10},
		{"return 10 return 9 return 8", 10},
		{"return 2 * 5 return 9", 10},
		{"9 return 2 * 5 return 9", 10},
		{"if (10 > 1) { return 10 }", 10},
		{"if (10 > 1) { if (10 > 1) { return 10 } return 1 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := evaluated.(*object.Integer)
		if !ok {
			t.Errorf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if integer.Value != tt.expected {
			t.Errorf("object has wrong value. got=%d, want=%d", integer.Value, tt.expected)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
		expectedLine    int
		expectedColumn  int
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
			1, 3,
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
			1, 3,
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
			1, 1,
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
			1, 6,
		},
		{
			"5 true + false 5",
			"unknown operator: BOOLEAN + BOOLEAN",
			1, 8,
		},
		{
			"if (10 > 1) { true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
			1, 20,
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
			4, 17,
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
			1, 9,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}

		if errObj.Line != tt.expectedLine {
			t.Errorf("wrong error line. expected=%d, got=%d",
				tt.expectedLine, errObj.Line)
		}

		if errObj.Column != tt.expectedColumn {
			t.Errorf("wrong error column. expected=%d, got=%d",
				tt.expectedColumn, errObj.Column)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5 a", 5},
		{"let a = 5 * 5 a", 25},
		{"let a = 5 let b = a b", 5},
		{"let a = 5 let b = a let c = a + b + 5 c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "\\x => x + 2"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		//{"fn identity(x) { x } identity(5)", 5},
		//{"let double = \\x => x * 2 double(5)", 10},
		{"fn add (x y) { x + y } add(5 5)", 10},
		{"let add = \\ x y => { x + y } add(5 + 5 add(5 5))", 20},
	}

	for i, tt := range tests {
		fmt.Printf("Test: %d\n", i)
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestPipeOperator(t *testing.T) {
	input := `
    fn identity (x) { x }
    fn add(x y) { x + y }
    let result = 5 |> identity() |> add(10);
    result
    `

	var expected int64 = 15

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, expected)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one" "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}
