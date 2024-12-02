// parser/parser.go

package parser

import (
	"fmt"
	"strconv"
	"unicode"

	"renelle/ast"
	"renelle/lexer"
	"renelle/token"
)

const (
	_ int = iota
	LOWEST
	OR          // or
	AND         // and
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	EXPONENT    // **
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	ACCESS      // object.property
)

var precedences = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NEQ:       EQUALS,
	token.ARRAY_EQ:  EQUALS,
	token.ARRAY_NEQ: EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.LTE:       LESSGREATER,
	token.GTE:       LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.CONCAT:    SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.MOD:       PRODUCT,
	token.POW:       EXPONENT,
	token.PIPE:      CALL,
	token.LPAREN:    CALL,
	token.FUNCCALL:  CALL,
	token.OR:        OR,
	token.AND:       AND,
	token.AT:        INDEX,
	token.DOTDOT:    INDEX,
	token.DOT:       ACCESS,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type ParseError struct {
	Message string
	Line    int
	Column  int
}

type Parser struct {
	l      *lexer.Lexer
	errors []ParseError

	curToken     token.Token
	peekToken    token.Token
	peekTokenTwo token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLitearl)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.INTERPOLATED, p.parseInterpolatedStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.COND, p.parseCondExpression)
	p.registerPrefix(token.CASE, p.parseCaseExpression)
	p.registerPrefix(token.BACKSLASH, p.parseFunctionLiteral)
	p.registerPrefix(token.FUNCCALL, p.parseCallExpression)
	p.registerPrefix(token.ATOM, p.parseAtom)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.ARRAY_EQ, p.parseInfixExpression)
	p.registerInfix(token.ARRAY_NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.CONCAT, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.POW, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AT, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parsePropertyAccessExpression)
	p.registerInfix(token.DOTDOT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.peekTokenTwo
	p.peekTokenTwo = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.MODULE:
		return p.parseModule()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	p.nextToken()

	switch p.curToken.Type {
	case token.IDENT:
		stmt.Left = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.LPAREN:
		stmt.Left = p.parseGroupedExpression()
	case token.LBRACKET:
		stmt.Left = p.parseArrayLiteral()
	case token.LBRACE:
		stmt.Left = p.parseMapLiteral()
	default:
		return nil
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	initialToken := p.curToken
	identifierValue := initialToken.Literal

	// If the identifier starts with an uppercase letter
	if unicode.IsUpper(rune(identifierValue[0])) {
		for {
			if !p.peekTokenIs(token.DOT) {
				break
			}
			// Look ahead one more token
			if p.peekTokenTwoIs(token.FUNCCALL) || (p.peekTokenTwoIs(token.IDENT) && unicode.IsLower(rune(p.peekTokenTwo.Literal[0]))) {
				break
			}
			p.nextToken() // Skip the dot
			identifierValue += "." + p.peekToken.Literal
			p.nextToken() // Move to the next identifier
		}
	}
	identifier := &ast.Identifier{Token: initialToken, Value: identifierValue}
	return identifier
}

func (p *Parser) parseAtom() ast.Expression {
	return &ast.AtomLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, ParseError{Message: msg, Line: p.curToken.Line, Column: p.curToken.Column})
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return lit
}

func (p *Parser) parseInterpolatedStringLiteral() ast.Expression {
	lit := &ast.InterpolatedStringLiteral{Token: p.curToken}
	lit.Segments = p.parseStringSegments(p.curToken.Literal)
	return lit
}

func (p *Parser) parseStringSegments(input string) []ast.Expression {
	segments := []ast.Expression{}
	var buffer string
	inExpression := false
	exprBuffer := ""

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '\\' && i+1 < len(input) && input[i+1] == '{' {
			// Handle escaped brace
			buffer += "{"
			i++ // Skip the next character
		} else if char == '{' {
			if inExpression {
				// Handle nested braces if necessary
				exprBuffer += string(char)
			} else {
				// Add the current buffer as a string segment
				if buffer != "" {
					literal := &ast.StringLiteral{Value: buffer}
					segments = append(segments, literal)
					buffer = ""
				}
				inExpression = true
			}
		} else if char == '}' {
			if inExpression {
				// Parse the expression inside the braces
				expr := p.parseExpressionFromString(exprBuffer)
				segments = append(segments, expr)
				exprBuffer = ""
				inExpression = false
			} else {
				buffer += string(char)
			}
		} else {
			if inExpression {
				exprBuffer += string(char)
			} else {
				buffer += string(char)
			}
		}
	}

	// Add any remaining buffer as a string segment
	if buffer != "" {
		literal := &ast.StringLiteral{Value: buffer}
		segments = append(segments, literal)
	}

	return segments
}

func (p *Parser) parseExpressionFromString(input string) ast.Expression {
	l := lexer.New(input, p.curToken.FileName)
	exprParser := New(l)
	expression := exprParser.parseExpression(LOWEST)
	return expression
}

func (p *Parser) parseFloatLitearl() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, ParseError{Message: msg, Line: p.curToken.Line, Column: p.curToken.Column})
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	p.nextToken()

	elements := []ast.Expression{}

	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		element := p.parseExpression(LOWEST)
		elements = append(elements, element)
		p.nextToken()
	}
	array.Elements = elements

	return array

}

func (p *Parser) parseMapLiteral() ast.Expression {
	if p.peekTokenIs(token.IDENT) && p.peekTokenTwoIs(token.WITH) {
		return p.parseMapUpdateLiteral()
	}

	mapLiteral := &ast.MapLiteral{Token: p.curToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if _, ok := key.(*ast.AtomLiteral); ok && !p.peekTokenIs(token.ASSIGN) {
			// If the key is an AtomLiteral and the next token is not ASSIGN, parse the key as an Atom
			p.nextToken()
			val := p.parseExpression(LOWEST)
			mapLiteral.Pairs[key] = val
		} else {
			// If the key is not an AtomLiteral or the next token is ASSIGN, parse the key as a String
			if !p.expectPeek(token.ASSIGN) {
				return nil
			}
			p.nextToken()
			value := p.parseExpression(LOWEST)
			mapLiteral.Pairs[key] = value
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return mapLiteral
}

func (p *Parser) parseMapUpdateLiteral() ast.Expression {
	mapUpdate := &ast.MapUpdateLiteral{Token: p.curToken}
	p.nextToken() // {
	mapUpdate.Left = p.parseIdentifier()
	p.nextToken() // with

	pairs := make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.ASSIGN) {
			return nil
		}
		p.nextToken()
		val := p.parseExpression(LOWEST)
		pairs[key] = val
	}
	p.nextToken()
	mapUpdate.Right = pairs
	return mapUpdate
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	initialToken := p.curToken
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	// if next token is a paren, this is a single expression grouping, else it's a tuple
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return expression
	}

	elements := []ast.Expression{expression}
	for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		elements = append(elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return &ast.TupleLiteral{Token: initialToken, Elements: elements}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	expression.Index = p.parseExpression(p.curPrecedence())

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseCondExpression() ast.Expression {
	expression := &ast.CondExpression{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		condition := p.parseExpression(LOWEST)
		expression.Conditions = append(expression.Conditions, condition)

		if !p.expectPeek(token.ARROW) {
			return nil
		}

		p.nextToken()

		if p.peekTokenIs(token.LBRACE) {
			consequence := p.parseBlockStatement()
			expression.Consequences = append(expression.Consequences, consequence)
		} else {
			consequence := p.parseExpression(LOWEST)
			expr := &ast.ExpressionStatement{Token: p.curToken, Expression: consequence}
			expression.Consequences = append(expression.Consequences, &ast.BlockStatement{Statements: []ast.Statement{expr}})

		}

	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return expression
}

func (p *Parser) parseCaseExpression() ast.Expression {
	expression := &ast.CaseExpression{Token: p.curToken}

	p.nextToken()

	expression.Test = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		condition := p.parseExpression(LOWEST)
		expression.Conditions = append(expression.Conditions, condition)

		if !p.expectPeek(token.ARROW) {
			return nil
		}

		p.nextToken()

		if p.curTokenIs(token.LBRACE) {
			consequence := p.parseBlockStatement()
			expression.Consequences = append(expression.Consequences, consequence)
		} else {
			consequence := p.parseExpression(LOWEST)
			expr := &ast.ExpressionStatement{Token: p.curToken, Expression: consequence}
			expression.Consequences = append(expression.Consequences, &ast.BlockStatement{Statements: []ast.Statement{expr}})

		}

	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return expression

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parsePropertyAccessExpression(left ast.Expression) ast.Expression {
	expression := &ast.PropertyAccessExpression{Token: p.curToken, Left: left}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken() // Consume the IDENT
		expression.Right = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else if p.peekTokenIs(token.FUNCCALL) {
		p.nextToken() // Consume the FUNCCALL
		expression.Right = p.parseCallExpression()
	} else {
		return nil
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	lit.Parameters = p.parseFunctionParameters()

	p.nextToken()

	if !p.curTokenIs(token.LBRACE) {
		lit.Body = &ast.BlockStatement{}
		exp := p.parseExpression(LOWEST)
		lit.Body.Statements = append(lit.Body.Statements, &ast.ExpressionStatement{Token: p.curToken, Expression: exp})
	} else {
		lit.Body = p.parseBlockStatement()
	}

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.ARROW) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for !p.peekTokenIs(token.ARROW) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	if p.peekToken.Type != token.IDENT && p.peekToken.Type != token.FUNCCALL {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionStatementParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionStatementParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression() ast.Expression {
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifier.Token.Type = token.IDENT

	exp := &ast.CallExpression{Token: p.curToken, Function: identifier}

	p.nextToken()
	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		exp.Arguments = []ast.Expression{}
		return exp
	}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	for {
		args = append(args, p.parseExpression(LOWEST))
		if p.peekTokenIs(token.RPAREN) || p.peekTokenIs(token.EOF) {
			break
		}
		p.nextToken() // Move to the next argument
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseModule() *ast.Module {
	tok := p.curToken

	// Expect the module name to be an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	moduleName := p.parseIdentifier().(*ast.Identifier)
	p.nextToken()

	// Parse the body of the module
	moduleBody := []ast.Statement{}
	for !p.peekTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			moduleBody = append(moduleBody, stmt)
		}
		p.nextToken()
	}

	return &ast.Module{Token: tok, Name: moduleName, Body: moduleBody}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekTokenTwoIs(t token.TokenType) bool {
	return p.peekTokenTwo.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d, col%d: expected %s, got %s instead", p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, ParseError{Message: msg, Line: p.peekToken.Line, Column: p.peekToken.Column})
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
