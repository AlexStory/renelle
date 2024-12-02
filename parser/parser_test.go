// parser/parser_test.go

package parser

import (
	"fmt"
	"renelle/ast"
	"renelle/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
let x  = 5
let y = 10
let foobar = 838383
`
	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {

		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %T\n", i, s)
		}
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	ident, _ := letStmt.Left.(*ast.Identifier)
	if ident.Value != name {
		t.Errorf("letStmt.Left.Value not '%s'. got=%s", name, ident.Value)
		return false
	}

	if ident.TokenLiteral() != name {
		t.Errorf("s.Ident.TokenLiteral not '%s'. got=%s", name, ident.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
    return 5
    return 10
    return 993322
    `

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	expectedValues := []int64{5, 10, 993322}

	for i, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

		// Check the ReturnValue
		literal, ok := returnStmt.ReturnValue.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("returnStmt.ReturnValue not *ast.IntegerLiteral. got=%T", returnStmt.ReturnValue)
			continue
		}
		if literal.Value != expectedValues[i] {
			t.Errorf("literal.Value not %d. got=%d", expectedValues[i], literal.Value)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg.Message)
	}
	t.FailNow()
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := "3.14"

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 3.14 {
		t.Errorf("literal.Value not %f. got=%f", 3.14, literal.Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input, "test")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not \"%d\". got=\"%d\"", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 1", 5, "+", 1},
		{"5 - 2", 5, "-", 2},
		{"5 * 3", 5, "*", 3},
		{"5 / 4", 5, "/", 4},
		{"5 > 5", 5, ">", 5},
		{"5 < 6", 5, "<", 6},
		{"5 == 7", 5, "==", 7},
		{"5 != 8", 5, "!=", 8},
		{"5 <= 9", 5, "<=", 9},
		{"5 >= 10", 5, ">=", 10},
		{"true and true", true, "and", true},
		{"true or false", true, "or", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input, "test")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			for i, s := range program.Statements {
				fmt.Printf("program.Statements[%d] = %T\n", i, s)
				fmt.Printf("%s\n", program.Statements[i].TokenLiteral())
			}
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{
			"a * ([1 2 3 4] @ b * c) * d",
			"((a * ([1 2 3 4] @ (b * c))) * d)",
		},
		{
			"add(a * b@2 b@1 2 * [1, 2]@1)",
			"add((a * (b @ 2)) (b @ 1) (2 * ([1 2] @ 1)))",
		}, {"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	input := "true"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[%d] is not ast.ExpressionStatement. got=%T", 0, program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, true)

	if !testBoolean(t, stmt.Expression, true) {
		return
	}
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil. got=%+v", exp.Alternative)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `
	\=> 5
    \ x => x * 2
    \ x y => { x + y }`

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedParamCount int
	}{
		{0},
		{1},
		{2},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		function, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not ast.ExpressionStatement. got=%T", i, stmt)
		}

		if len(function.Expression.(*ast.FunctionLiteral).Parameters) != tt.expectedParamCount {
			t.Fatalf("function literal parameters wrong. want %d, got=%d\n", tt.expectedParamCount, len(function.Expression.(*ast.FunctionLiteral).Parameters))
		}

		body := function.Expression.(*ast.FunctionLiteral).Body
		if len(body.Statements) == 0 {
			t.Fatalf("block statement is empty")
		}
	}
}

func TestFunctionStatement(t *testing.T) {
	input := `fn add(x y) { x + y }`

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%q", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(stmt.Parameters))
	}

	if stmt.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", stmt.Parameters[0].String())
	}

	if stmt.Parameters[1].String() != "y" {
		t.Fatalf("parameter is not 'y'. got=%q", stmt.Parameters[1].String())
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("body.Statements has not 1 statements. got=%d", len(stmt.Body.Statements))
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1 (2 * 3) (4 + 5))"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		for i, a := range exp.Arguments {
			fmt.Printf("exp.Arguments[%d] = %s\n", i, a.String())
		}
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestAtomExpressionParsing(t *testing.T) {
	input := ":ok"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	atom, ok := stmt.Expression.(*ast.AtomLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.AtomExpression. got=%T",
			stmt.Expression)
	}

	if atom.Value != "ok" {
		t.Errorf("atom.Value not %s. got=%s", "ok", atom.Value)
	}
}

func TestParserIgnoresComments(t *testing.T) {
	input := `# comment before function
fn add(x y) { # comment at end of line
    x + y # comment at end of line
} # comment at end of line
# comment on its own line
fn subtract(x y) { # comment at end of line
    x - y # comment at end of line
} # comment at end of line`

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			2, len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"add"},
		{"subtract"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testFunctionStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testFunctionStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "fn" {
		t.Errorf("s.TokenLiteral not 'fn'. got=%q", s.TokenLiteral())
		return false
	}

	fnStmt, ok := s.(*ast.FunctionStatement)
	if !ok {
		t.Errorf("s not *ast.FunctionStatement. got=%T", s)
		return false
	}

	if fnStmt.Name.Value != name {
		t.Errorf("function name not '%s'. got=%s", name, fnStmt.Name.Value)
		return false
	}

	if fnStmt.Name.TokenLiteral() != name {
		t.Errorf("fnStmt.Name not '%s'. got=%s", name, fnStmt.Name)
		return false
	}

	return true
}

func TestLetReturnStatement(t *testing.T) {
	input := "let a = 5 a"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	if _, ok := program.Statements[0].(*ast.LetStatement); !ok {
		t.Errorf("program.Statements[0] not *ast.LetStatement. got=%T", program.Statements[0])
	}

	if _, ok := program.Statements[1].(*ast.ExpressionStatement); !ok {
		t.Errorf("program.Statements[1] not *ast.ExpressionStatement. got=%T", program.Statements[1])
	}
}

func TestMultipleLetStatements(t *testing.T) {
	input := "let a = 5 let b = a let c = a + b + 5 c"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain 4 statements. got=%d",
			len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.LetStatement. got=%T", program.Statements[0])
	}

	_, ok = program.Statements[1].(*ast.LetStatement)
	if !ok {
		t.Errorf("program.Statements[1] not *ast.LetStatement. got=%T", program.Statements[1])
	}

	_, ok = program.Statements[2].(*ast.LetStatement)
	if !ok {
		t.Errorf("program.Statements[2] not *ast.LetStatement. got=%T", program.Statements[2])
	}

	_, ok = program.Statements[3].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[3] not *ast.ExpressionStatement. got=%T", program.Statements[3])
	}
}

func TestFunctionCallParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectedArgs int
	}{
		{"identity()", 0},
		{"identity(arg)", 1},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		call, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
		}

		if len(call.Arguments) != tt.expectedArgs {
			t.Errorf("wrong number of arguments. expected=%d, got=%d", tt.expectedArgs, len(call.Arguments))
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}
func TestInterpolatedStringLiteral_NoExpressions(t *testing.T) {
	input := `$"hello world"`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.InterpolatedStringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.InterpolatedStringLiteral. got=%T", stmt.Expression)
	}

	if len(literal.Segments) != 1 {
		t.Fatalf("wrong number of segments. expected=1, got=%d", len(literal.Segments))
	}

	strSegment, ok := literal.Segments[0].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("segment not *ast.StringLiteral. got=%T", literal.Segments[0])
	}

	if strSegment.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", strSegment.Value)
	}
}

func TestInterpolatedStringLiteral_OneExpression(t *testing.T) {
	input := `$"hello {name}"`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.InterpolatedStringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.InterpolatedStringLiteral. got=%T", stmt.Expression)
	}

	if len(literal.Segments) != 2 {
		t.Fatalf("wrong number of segments. expected=2, got=%d", len(literal.Segments))
	}

	strSegment, ok := literal.Segments[0].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("segment not *ast.StringLiteral. got=%T", literal.Segments[0])
	}

	if strSegment.Value != "hello " {
		t.Errorf("literal.Value not %q. got=%q", "hello ", strSegment.Value)
	}

	exprSegment, ok := literal.Segments[1].(*ast.Identifier)
	if !ok {
		t.Fatalf("segment not *ast.Identifier. got=%T", literal.Segments[1])
	}

	if exprSegment.Value != "name" {
		t.Errorf("identifier.Value not %q. got=%q", "name", exprSegment.Value)
	}
}

func TestInterpolatedStringLiteral_MultipleExpressions(t *testing.T) {
	input := `$"hello {name}, you have {count} new messages"`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.InterpolatedStringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.InterpolatedStringLiteral. got=%T", stmt.Expression)
	}

	if len(literal.Segments) != 5 {
		t.Fatalf("wrong number of segments. expected=5, got=%d", len(literal.Segments))
	}

	strSegment1, ok := literal.Segments[0].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("segment not *ast.StringLiteral. got=%T", literal.Segments[0])
	}

	if strSegment1.Value != "hello " {
		t.Errorf("literal.Value not %q. got=%q", "hello ", strSegment1.Value)
	}

	exprSegment1, ok := literal.Segments[1].(*ast.Identifier)
	if !ok {
		t.Fatalf("segment not *ast.Identifier. got=%T", literal.Segments[1])
	}

	if exprSegment1.Value != "name" {
		t.Errorf("identifier.Value not %q. got=%q", "name", exprSegment1.Value)
	}

	strSegment2, ok := literal.Segments[2].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("segment not *ast.StringLiteral. got=%T", literal.Segments[2])
	}

	if strSegment2.Value != ", you have " {
		t.Errorf("literal.Value not %q. got=%q", ", you have ", strSegment2.Value)
	}

	exprSegment2, ok := literal.Segments[3].(*ast.Identifier)
	if !ok {
		t.Fatalf("segment not *ast.Identifier. got=%T", literal.Segments[3])
	}

	if exprSegment2.Value != "count" {
		t.Errorf("identifier.Value not %q. got=%q", "count", exprSegment2.Value)
	}

	strSegment3, ok := literal.Segments[4].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("segment not *ast.StringLiteral. got=%T", literal.Segments[4])
	}

	if strSegment3.Value != " new messages" {
		t.Errorf("literal.Value not %q. got=%q", " new messages", strSegment3.Value)
	}
}
func TestParseStringSegments(t *testing.T) {
	tests := []struct {
		input    string
		expected []ast.Expression
	}{
		{
			input: `hello world`,
			expected: []ast.Expression{
				&ast.StringLiteral{Value: "hello world"},
			},
		},
		{
			input: `hello {name}`,
			expected: []ast.Expression{
				&ast.StringLiteral{Value: "hello "},
				&ast.Identifier{Value: "name"},
			},
		},
		{
			input: `hello {name}, you have {count} new messages`,
			expected: []ast.Expression{
				&ast.StringLiteral{Value: "hello "},
				&ast.Identifier{Value: "name"},
				&ast.StringLiteral{Value: ", you have "},
				&ast.Identifier{Value: "count"},
				&ast.StringLiteral{Value: " new messages"},
			},
		},
		{
			input: `hello \{name}`,
			expected: []ast.Expression{
				&ast.StringLiteral{Value: "hello {name}"},
			},
		},
	}

	for _, tt := range tests {
		p := &Parser{}
		segments := p.parseStringSegments(tt.input)

		if len(segments) != len(tt.expected) {
			t.Fatalf("wrong number of segments. expected=%d, got=%d", len(tt.expected), len(segments))
		}

		for i, segment := range segments {
			switch expected := tt.expected[i].(type) {
			case *ast.StringLiteral:
				actual, ok := segment.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("segment not *ast.StringLiteral. got=%T", segment)
				}
				if actual.Value != expected.Value {
					t.Errorf("literal.Value not %q. got=%q", expected.Value, actual.Value)
				}
			case *ast.Identifier:
				actual, ok := segment.(*ast.Identifier)
				if !ok {
					t.Fatalf("segment not *ast.Identifier. got=%T", segment)
				}
				if actual.Value != expected.Value {
					t.Errorf("identifier.Value not %q. got=%q", expected.Value, actual.Value)
				}
			}
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1 2 * 2 3 + 3]"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray @ 1 + 1"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingTuples(t *testing.T) {
	input := "(1 2 3 4 5)"

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	tuple, ok := stmt.Expression.(*ast.TupleLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.TupleLiteral. got=%T", stmt.Expression)
	}

	if len(tuple.Elements) != 5 {
		t.Fatalf("tuple does not contain %d elements. got=%d",
			5, len(tuple.Elements))
	}

	for i, elem := range tuple.Elements {
		integer, ok := elem.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("tuple.Elements[%d] is not ast.IntegerLiteral. got=%T", i, elem)
		}

		if integer.Value != int64(i+1) {
			t.Errorf("integer.Value is not %d. got=%d", i+1, integer.Value)
		}
	}
}
func TestParsingMapLiteralsStringKeys(t *testing.T) {
	input := `{"one" = 1, "two" = 2, "three" = 3}`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hmap, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hmap.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hmap.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hmap.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one" = 0 + 1, "two" = 10 - 8, "three" = 15 / 5}`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func TestParsingMapLiteralsAtomKeys(t *testing.T) {
	input := `{name: "Hayley", age: 8}`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hmap, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.MapLiteral. got=%T", stmt.Expression)
	}

	if len(hmap.Pairs) != 2 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hmap.Pairs))
	}

	expected := map[string]interface{}{
		"name": "Hayley",
		"age":  8,
	}

	for key, value := range hmap.Pairs {
		literal, ok := key.(*ast.AtomLiteral)
		if !ok {
			t.Errorf("key is not ast.AtomLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		switch expectedValue := expectedValue.(type) {
		case int:
			testIntegerLiteral(t, value, int64(expectedValue))
		case string:
			testStringLiteral(t, value, expectedValue)
		default:
			t.Errorf("Type not handled. got=%T", expectedValue)
		}
	}
}

func testStringLiteral(t *testing.T, il ast.Expression, value string) {
	strlit, ok := il.(*ast.StringLiteral)
	if !ok {
		t.Errorf("il not *ast.StringLiteral. got=%T", il)
		return
	}

	if strlit.Value != value {
		t.Errorf("strlit.Value not %q. got=%q", value, strlit.Value)
	}
}

func TestPropertyAccessExpression(t *testing.T) {
	input := `myMap.property`

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	access, ok := stmt.Expression.(*ast.PropertyAccessExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.PropertyAccessExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, access.Left, "myMap") {
		return
	}

	if !testIdentifier(t, access.Right, "property") {
		return
	}
}

func TestCondExpression(t *testing.T) {
	input := `
    cond {
        x > 10 => "greater than 10"
        x > 5 => "greater than 5"
        true => "default"
    }
    `

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	cond, ok := stmt.Expression.(*ast.CondExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CondExpression. got=%T", stmt.Expression)
	}

	if len(cond.Conditions) != 3 {
		t.Fatalf("conditions does not contain 3 conditions. got=%d", len(cond.Conditions))
	}

	if len(cond.Consequences) != 3 {
		t.Fatalf("consequences does not contain 3 consequences. got=%d", len(cond.Consequences))
	}
}

func TestCaseExpression(t *testing.T) {
	input := `
    case run_function() {
        (:ok, data) => "do good thing"
        (:err, error) => "handle error"
    }
    `

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	caseExpr, ok := stmt.Expression.(*ast.CaseExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CaseExpression. got=%T", stmt.Expression)
	}

	if len(caseExpr.Conditions) != 2 {
		t.Fatalf("conditions does not contain 2 conditions. got=%d", len(caseExpr.Conditions))
	}

	if len(caseExpr.Consequences) != 2 {
		t.Fatalf("consequences does not contain 2 consequences. got=%d", len(caseExpr.Consequences))
	}
}
func TestMultilineCaseExpression(t *testing.T) {
	input := `
    case x {
        n => {
            let x = n + 2
            print(x)
        }
    }
    `

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	caseExpr, ok := stmt.Expression.(*ast.CaseExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CaseExpression. got=%T", stmt.Expression)
	}

	if len(caseExpr.Conditions) != 1 {
		t.Fatalf("conditions does not contain 1 condition. got=%d", len(caseExpr.Conditions))
	}

	if len(caseExpr.Consequences) != 1 {
		t.Fatalf("consequences does not contain 1 consequence. got=%d", len(caseExpr.Consequences))
	}

	condition := caseExpr.Conditions[0]
	if condition.String() != "n" {
		t.Fatalf("condition is not 'n'. got=%s", condition.String())
	}

	consequence := caseExpr.Consequences[0]

	if len(consequence.Statements) != 2 {
		t.Fatalf("consequence does not contain 2 statements. got=%d", len(consequence.Statements))
	}

	letStmt, ok := consequence.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not ast.LetStatement. got=%T", consequence.Statements[0])
	}

	if letStmt.Left.TokenLiteral() != "x" {
		t.Fatalf("letStmt.Name.Value is not 'x'. got=%s", letStmt.TokenLiteral())
	}

	if letStmt.Value.String() != "(n + 2)" {
		t.Fatalf("letStmt.Value is not '(n + 2)'. got=%s", letStmt.Value.String())
	}

	exprStmt, ok := consequence.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[1] is not ast.ExpressionStatement. got=%T", consequence.Statements[1])
	}

	if exprStmt.Expression.String() != "print(x)" {
		t.Fatalf("exprStmt.Expression is not 'print(x)'. got=%s", exprStmt.Expression.String())
	}
}

func TestParseModuleIdentifier(t *testing.T) {
	input := "MyApp.MyModule.MySubmodule"
	l := lexer.New(input, "test")
	p := New(l)

	identifier := p.parseIdentifier().(*ast.Identifier)
	if identifier.Value != input {
		t.Errorf("identifier.Value not '%s'. got=%s", input, identifier.Value)
	}

	input = "MyApp.MyModule.func"
	l = lexer.New(input, "test")
	p = New(l)

	identifier = p.parseIdentifier().(*ast.Identifier)
	if identifier.Value != "MyApp.MyModule" {
		t.Errorf("identifier.Value not 'MyApp.MyModule'. got=%s", identifier.Value)
	}
}

func TestParseModule(t *testing.T) {
	input := `
    module MyApp.Module
    fn x() {
    }
    `

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	moduleStmt, ok := program.Statements[0].(*ast.Module)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Module. got=%T",
			program.Statements[0])
	}

	if moduleStmt.Name.Value != "MyApp.Module" {
		t.Errorf("moduleStmt.Name.Value not 'MyApp.Module'. got=%q", moduleStmt.Name.Value)
	}

	if len(moduleStmt.Body) != 1 {
		for i, s := range moduleStmt.Body {
			fmt.Printf("moduleStmt.Body[%d] = %s\n", i, s.String())
		}
		t.Fatalf("moduleStmt.Body does not contain 1 statements. got=%d",
			len(moduleStmt.Body))
	}
}

func TestModuleAccessExpression(t *testing.T) {
	input := `Mod.func()`

	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	access, ok := stmt.Expression.(*ast.PropertyAccessExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.PropertyAccessExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, access.Left, "Mod") {
		return
	}

	call, ok := access.Right.(*ast.CallExpression)
	if !ok {
		t.Fatalf("access.Right is not ast.CallExpression. got=%T",
			access.Right)
	}

	if !testIdentifier(t, call.Function, "func") {
		return
	}
}

func TestMapUpdateLiteral(t *testing.T) {
	input := `{ identifier with key = "value", anotherKey = "anotherValue" }`
	l := lexer.New(input, "test")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		for i, s := range program.Statements {
			fmt.Printf("program.Statements[%d] = %s\n", i, s.String())
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	mapUpdate, ok := stmt.Expression.(*ast.MapUpdateLiteral)
	if !ok {
		t.Fatalf("exp not *ast.MapUpdateLiteral. got=%T", stmt.Expression)
	}

	if mapUpdate.Left.String() != "identifier" {
		t.Fatalf("mapUpdate.Left not 'identifier'. got=%q", mapUpdate.Left.String())
	}

	if len(mapUpdate.Right) != 2 {
		t.Fatalf("mapUpdate does not contain %d pairs. got=%d",
			2, len(mapUpdate.Right))
	}

	testMapUpdateLiteral(t, mapUpdate)
}

func testMapUpdateLiteral(t *testing.T, ml *ast.MapUpdateLiteral) {
	if ml.TokenLiteral() != "{" {
		t.Errorf("ml.TokenLiteral not '{'. got=%q", ml.TokenLiteral())
	}

	if len(ml.Right) != 2 {
		t.Errorf("ml.Pairs has wrong length. got=%d", len(ml.Right))
	}

	tests := map[string]string{
		"key":        "value",
		"anotherKey": "anotherValue",
	}

	for key, value := range ml.Right {
		literal := key.(*ast.Identifier).Value
		if tests[literal] != value.(*ast.StringLiteral).Value {
			t.Errorf("ml.Pairs[%q] not '%s'. got=%q", literal, tests[literal], value.(*ast.StringLiteral).Value)
		}
	}
}
