// evaluator/evaluator_test.go

package evaluator

import (
	"fmt"
	"renelle/constants"
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
	l := lexer.New(input, "test")
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	md := object.NewEvalContext()

	return Eval(program, env, md)
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
		{":foo", "foo"},
		{":bar", "bar"},
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
	if obj != constants.NIL {
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

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3]@0",
			1,
		},
		{
			"[1, 2, 3]@1",
			2,
		},
		{
			"[1, 2, 3]@2",
			3,
		},
		{
			"let i = 0 [1]@i",
			1,
		},
		{
			"[1, 2, 3]@1 + 1",
			3,
		},
		{
			"let myArray = [1, 2, 3] myArray@2",
			3,
		},
		{
			"let myArray = [1 2 3] (myArray@0) + (myArray@1) + (myArray@2)",
			6,
		},
		{
			"let myArray = [1, 2, 3] let i = myArray@0 myArray@i",
			2,
		},
		{
			"[1, 2, 3]@3",
			nil,
		},
		{
			"[1, 2, 3]@-1",
			nil,
		},
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

func TestEvalTuples(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(1 2 3)", "(1 2 3)"},
		{"(1 + 2 3 * 4 5)", "(3 12 5)"},
		{"(1 (2 3) 4)", "(1 (2 3) 4)"},
		{"let a = 5 (a a + 1)", "(5 6)"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		tuple, ok := evaluated.(*object.Tuple)
		if !ok {
			t.Fatalf("object is not Tuple. got=%T (%+v)", evaluated, evaluated)
		}

		if tuple.Inspect() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, tuple.Inspect())
		}
	}
}

func TestArrayDestructuring(t *testing.T) {
	// Test successful destructuring
	result := testEval("let [x, y] = [1, 2]; x + y")
	if resultIns, ok := result.(*object.Integer); ok {
		if resultIns.Value != 3 {
			t.Errorf("Expected 3, got %v", resultIns.Value)
		}
	} else {
		t.Errorf("Expected no error, got %v", result)
	}

	// Test destructuring with discards
	result = testEval("let [_, y] = [1, 2]; y")
	if resultIns, ok := result.(*object.Integer); ok {
		if resultIns.Value != 2 {
			t.Errorf("Expected 2, got %v", resultIns.Value)
		}
	} else {
		t.Errorf("Expected no error, got %v", result)
	}

	// Test mismatched length
	result = testEval("let [x, y] = [1]; x + y")
	if _, ok := result.(*object.Error); !ok {
		t.Errorf("Expected error, got nil")
	}

	// Test mismatched literal
	result = testEval("let [1, y] = [2, 2]; y")
	if _, ok := result.(*object.Error); !ok {
		t.Errorf("Expected error, got nil")
	}
}

func TestTupleDestructuring(t *testing.T) {
	// Test successful destructuring
	result := testEval("let (x, y) = (1, 2); x + y")
	if resultIns, ok := result.(*object.Integer); ok {
		if resultIns.Value != 3 {
			t.Errorf("Expected 3, got %v", resultIns.Value)
		}
	} else {
		t.Errorf("Expected no error, got %v", result)
	}

	// Test destructuring with discards
	result = testEval("let (_, y) = (1, 2); y")
	if resultIns, ok := result.(*object.Integer); ok {
		if resultIns.Value != 2 {
			t.Errorf("Expected 2, got %v", resultIns.Value)
		}
	} else {
		t.Errorf("Expected no error, got %v", result)
	}

	// Test mismatched length
	result = testEval("let (x, y) = (1); x + y")
	if _, ok := result.(*object.Error); !ok {
		t.Errorf("Expected error, got nil")
	}

	// Test mismatched literal
	result = testEval("let (1, y) = (2, 2); y")
	if _, ok := result.(*object.Error); !ok {
		t.Errorf("Expected error, got nil")
	}
}

func TestEvalMapLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected map[interface{}]interface{}
	}{
		{
			input: `let cat = {name: "hayley", age: 8}`,
			expected: map[interface{}]interface{}{
				&object.Atom{Value: "name"}: "hayley",
				&object.Atom{Value: "age"}:  8,
			},
		},
		{
			input: `let dog = {"name" = "goldie", "age" = 8}`,
			expected: map[interface{}]interface{}{
				&object.String{Value: "name"}: "goldie",
				&object.String{Value: "age"}:  8,
			},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Map)
		if !ok {
			t.Fatalf("object is not Map. got=%T (%+v)", evaluated, evaluated)
		}

		if len(result.Store.Buckets) != len(tt.expected) {
			t.Fatalf("Map has wrong num of pairs. got=%d, want=%d",
				len(result.Store.Buckets), len(tt.expected))
		}

		for expectedKey, expectedValue := range tt.expected {
			value, ok := result.Get(expectedKey.(object.Object))
			if !ok {
				t.Fatalf("no value for given key in Map")
			}

			switch expected := expectedValue.(type) {
			case string:
				str, ok := value.(*object.String)
				if !ok {
					t.Errorf("value is not *object.String. got=%T (%+v)", value, value)
					continue
				}

				if str.Value != expected {
					t.Errorf("value is not %q. got=%q", expected, str.Value)
				}
			case int:
				integer, ok := value.(*object.Integer)
				if !ok {
					t.Errorf("value is not *object.Integer. got=%T (%+v)", value, value)
					continue
				}

				if integer.Value != int64(expected) {
					t.Errorf("value is not %d. got=%d", expected, integer.Value)
				}
			}
		}
	}
}
func TestEvalMapIndexOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			input:    `let cat = {name: "hayley", age: 8}; cat@:name`,
			expected: "hayley",
		},
		{
			input:    `let dog = {"name" = "goldie", "age" = 8}; dog@"name"`,
			expected: "goldie",
		},
		{
			input:    `let cat = {name: "hayley", age: 8}; cat@:age`,
			expected: 8,
		},
		{
			input:    `let dog = {"name" = "goldie", "age" = 8}; dog@"age"`,
			expected: 8,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if str.Value != expected {
				t.Errorf("String object has wrong value. got=%q, want=%q", str.Value, expected)
			}
		case int:
			integer, ok := evaluated.(*object.Integer)
			if !ok {
				t.Errorf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if integer.Value != int64(expected) {
				t.Errorf("Integer object has wrong value. got=%d, want=%d", integer.Value, expected)
			}
		}
	}
}

func TestPropertyAccessExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let m = {property: 5} m.property", 5},
		{"let m = {property: 5, anotherProperty: 10}; m.anotherProperty", 10},
		{"let m = {property: 5}; m.unknownProperty", constants.NIL},
		{"let m = 5; m.property", "property access not supported: INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		default:
			testNilObject(t, evaluated)
		}
	}
}

func TestCondExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"cond { true => 1 }", 1},
		{"cond { false => 1 true => 2 }", 2},
		{"cond { false => 1 false => 2 }", nil},
		{"cond { 1 > 2 => 1 2 > 1 => 2 }", 2},
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

func TestEvalCaseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let x = 1; case x { 1 => 2, _ => 3 }",
			2,
		},
		{
			"let x = 2; case x { 1 => 2, _ => 3 }",
			3,
		},
		{
			"let x = (1, 2); case x { (1, 2) => 3, _ => 4 }",
			3,
		},
		{
			"let x = (2, 3); case x { (1, 2) => 3, _ => 4 }",
			4,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func TestEvalModule(t *testing.T) {
	input := `
    module TestModule 

    let x = 10;
    let y = 20;
    x + y;
    
    `

	evaluated := testEval(input)
	module, ok := evaluated.(*object.Module)
	if !ok {
		t.Fatalf("object is not Module. got=%T (%+v)", evaluated, evaluated)
	}

	if module.Name != "TestModule" {
		t.Errorf("module has wrong name. got=%q", module.Name)
	}

	val, ok := module.Environment.Get("x")
	if !ok {
		t.Errorf("variable 'x' not found in module")
	}

	testIntegerObject(t, val, 10)

	val, ok = module.Environment.Get("y")
	if !ok {
		t.Errorf("variable 'y' not found in module")
	}

	testIntegerObject(t, val, 20)
}

func TestMapUpdateLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let m = {key: "value"}; let m2 = {m with :newKey = "newValue"}; m2.newKey`, "newValue"},
		{`let m = {key: "value"}; let m2 = {m with :key = "newValue"}; m2.key`, "newValue"},
		{`let m = {key: "value"}; let m2 = {m with :unknownKey = "newValue"}; m2.unknownKey`, "newValue"},
		{`let m = 5; let m2 = {m with :key = "value"}; m2.key`, "not a map: INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			strObj, ok := evaluated.(*object.String)
			if !ok {
				errObj, ok := evaluated.(*object.Error)
				if ok {
					if errObj.Message != expected {
						t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
					}
				} else {
					t.Errorf("object is not String or Error. got=%T (%+v)", evaluated, evaluated)
				}
				continue
			}
			if strObj.Value != expected {
				t.Errorf("wrong string value. expected=%q, got=%q", expected, strObj.Value)
			}
		default:
			testNilObject(t, evaluated)
		}
	}
}
