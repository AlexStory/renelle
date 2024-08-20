// lexer/lexer.go

package lexer

import (
	"bytes"
	"renelle/token"
	"strings"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	for {
		l.skipWhitespace()
		l.skipComments()
		if l.ch != ' ' && l.ch != '#' {
			break
		}
	}

	switch l.ch {
	case '=':
		if l.getNextChar() == '=' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Column: col}
		} else if l.getNextChar() == '>' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.ARROW, Literal: literal, Line: l.line, Column: col}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch, l.line, l.column)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.line, l.column)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.line, l.column)
	case '\\':
		tok = newToken(token.BACKSLASH, l.ch, l.line, l.column)
	case '%':
		tok = newToken(token.MOD, l.ch, l.line, l.column)
	case '*':
		if l.getNextChar() == '*' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.POW, Literal: literal, Line: l.line, Column: col}
		} else {
			tok = newToken(token.ASTERISK, l.ch, l.line, l.column)
		}
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.column)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line, l.column)
	case '@':
		tok = newToken(token.AT, l.ch, l.line, l.column)
	case '<':
		if l.getNextChar() == '=' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: l.line, Column: col}
		} else {
			tok = newToken(token.LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.getNextChar() == '=' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: l.line, Column: col}
		} else {
			tok = newToken(token.GT, l.ch, l.line, l.column)
		}
	case '!':
		if l.getNextChar() == '=' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NEQ, Literal: literal, Line: l.line, Column: col}
		} else {
			tok = newToken(token.BANG, l.ch, l.line, l.column)
		}
	case '|':
		if l.getNextChar() == '>' {
			ch := l.ch
			col := l.column
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PIPE, Literal: literal, Line: l.line, Column: col}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	case ':':
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = l.readAtom()
		tok.Type = token.ATOM
		return tok
	case '.':
		tok = newToken(token.DOT, l.ch, l.line, l.column)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column

	case '"':
		tok.Line = l.line
		tok.Column = l.column
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '$':
		tok.Line = l.line
		tok.Column = l.column
		l.readChar() // skip the $
		if l.ch == '"' {
			tok.Type = token.INTERPOLATED
			tok.Literal = l.readString()
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}

	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			switch l.ch {
			case '(':
				tok.Type = token.FUNCCALL
			case ':':
				tok.Type = token.ATOM
				l.readChar()
			default:
				tok.Type = token.LookupIdent(tok.Literal)
			}

			return tok
		} else if isDigit(l.ch, l.getPrevChar(), l.getNextChar()) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	isNewLine := l.ch == '\n'
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL"
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
	// Update line and column
	if isNewLine {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

}

func (l *Lexer) skipComments() {
	if l.ch == '#' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}

		if l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' || l.ch == ',' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte, line int, column int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?' || ch == '!'
}

func isDigit(ch byte, prevChar byte, nextChar byte) bool {
	if ch == '.' && isDigit(prevChar, 0, 0) && isDigit(nextChar, 0, 0) {
		return true
	}
	if ch == '_' {
		// If the previous character is not a digit, then '_' is the first character of the number
		return prevChar >= '0' && prevChar <= '9'
	}
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch, l.getPrevChar(), l.getNextChar()) {
		l.readChar()
	}
	return strings.ReplaceAll(l.input[position:l.position], "_", "")
}

func (l *Lexer) readAtom() string {
	position := l.position + 1
	l.readChar()

	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) getPrevChar() byte {
	if l.position > 0 {
		return l.input[l.position-1]
	}
	return 0
}

func (l *Lexer) getNextChar() byte {
	if l.position+1 < len(l.input) {
		return l.input[l.position+1]
	}
	return 0
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch, 0, 0) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var out bytes.Buffer
	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out.WriteString("\n")
			case 't':
				out.WriteString("\t")
			case 'r':
				out.WriteString("\r")
			case '"':
				out.WriteString("\"")
			case '\\':
				out.WriteString("\\")
			default:
				out.WriteString("\\" + string(l.ch))
			}

		} else if l.ch == '"' || l.ch == 0 {
			break
		} else {
			out.WriteByte(l.ch)
		}
	}
	return out.String()
}
