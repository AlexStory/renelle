// token/token.go

package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NEWLINE = "NEWLINE"

	// Identifiers + literals

	TRUE     = "TRUE"
	FALSE    = "FALSE"
	LET      = "LET"
	IDENT    = "IDENT"
	FUNCTION = "FUNCTION"
	ATOM     = "ATOM"
	FUNCCALL = "FUNCCALL"

	IF        = "IF"
	ELSE      = "ELSE"
	COND      = "COND"
	CASE      = "CASE"
	RETURN    = "RETURN"
	INT       = "INT"
	FLOAT     = "FLOAT"
	STRING    = "STRING"
	AND       = "AND"
	OR        = "OR"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	LT        = "<"
	GT        = ">"
	ASTERISK  = "*"
	MOD       = "%"
	SLASH     = "/"
	AT        = "@"
	BACKSLASH = "\\"
	ARROW     = "=>"
	EQ        = "=="
	LTE       = "<="
	GTE       = ">="
	NEQ       = "!="
	POW       = "**"
	PIPE      = "|>"
	DOT       = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var TokenMap = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"cond":   COND,
	"case":   CASE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"and":    AND,
	"or":     OR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := TokenMap[ident]; ok {
		return tok
	}
	return IDENT
}
