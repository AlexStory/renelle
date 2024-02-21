// lexer/lexer_test.go

package lexer

import (
	"fmt"
	"renelle/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5
let ten = 10
fn add (x y) {
	x + y
}
let result = add(x y)
!5<>10
if 5 < 10 {
	return true
} else {
	return false
}
4 <= 5 >= 3
5 == 5
5 != 9
2 ** 3
*/
2 |> add(3)
and or
:ok
[]%
3.14
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LET, "let", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		{token.LET, "let", 2, 1},
		{token.IDENT, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "10", 2, 11},
		{token.FUNCTION, "fn", 3, 1},
		{token.IDENT, "add", 3, 4},
		{token.LPAREN, "(", 3, 8},
		{token.IDENT, "x", 3, 9},
		{token.IDENT, "y", 3, 11},
		{token.RPAREN, ")", 3, 12},
		{token.LBRACE, "{", 3, 14},
		{token.IDENT, "x", 4, 2},
		{token.PLUS, "+", 4, 4},
		{token.IDENT, "y", 4, 6},
		{token.RBRACE, "}", 5, 1},
		{token.LET, "let", 6, 1},
		{token.IDENT, "result", 6, 5},
		{token.ASSIGN, "=", 6, 12},
		{token.IDENT, "add", 6, 14},
		{token.LPAREN, "(", 6, 17},
		{token.IDENT, "x", 6, 18},
		{token.IDENT, "y", 6, 20},
		{token.RPAREN, ")", 6, 21},
		{token.BANG, "!", 7, 1},
		{token.INT, "5", 7, 2},
		{token.LT, "<", 7, 3},
		{token.GT, ">", 7, 4},
		{token.INT, "10", 7, 5},
		{token.IF, "if", 8, 1},
		{token.INT, "5", 8, 4},
		{token.LT, "<", 8, 6},
		{token.INT, "10", 8, 8},
		{token.LBRACE, "{", 8, 11},
		{token.RETURN, "return", 9, 2},
		{token.TRUE, "true", 9, 9},
		{token.RBRACE, "}", 10, 1},
		{token.ELSE, "else", 10, 3},
		{token.LBRACE, "{", 10, 8},
		{token.RETURN, "return", 11, 2},
		{token.FALSE, "false", 11, 9},
		{token.RBRACE, "}", 12, 1},
		{token.INT, "4", 13, 1},
		{token.LTE, "<=", 13, 3},
		{token.INT, "5", 13, 6},
		{token.GTE, ">=", 13, 8},
		{token.INT, "3", 13, 11},
		{token.INT, "5", 14, 1},
		{token.EQ, "==", 14, 3},
		{token.INT, "5", 14, 6},
		{token.INT, "5", 15, 1},
		{token.NEQ, "!=", 15, 3},
		{token.INT, "9", 15, 6},
		{token.INT, "2", 16, 1},
		{token.POW, "**", 16, 3},
		{token.INT, "3", 16, 6},
		{token.ASTERISK, "*", 17, 1},
		{token.SLASH, "/", 17, 2},
		{token.INT, "2", 18, 1},
		{token.PIPE, "|>", 18, 3},
		{token.IDENT, "add", 18, 6},
		{token.LPAREN, "(", 18, 9},
		{token.INT, "3", 18, 10},
		{token.RPAREN, ")", 18, 11},
		{token.AND, "and", 19, 1},
		{token.OR, "or", 19, 5},
		{token.ATOM, ":ok", 20, 1},
		{token.LBRACKET, "[", 21, 1},
		{token.RBRACKET, "]", 21, 2},
		{token.MOD, "%", 21, 3},
		{token.FLOAT, "3.14", 22, 1},
		{token.EOF, "", 23, 1},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Printf("Token: %v", tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func Test0Floats(t *testing.T) {
	input := `0.14`

	l := New(input)
	tok := l.NextToken()
	fmt.Printf("Token: %v", tok)
	if tok.Type != token.FLOAT {
		t.Fatalf("expected float, got %v", tok.Type)
	}

	if tok.Literal != "0.14" {
		t.Fatalf("expected 0.14, got %v", tok.Literal)
	}
}
