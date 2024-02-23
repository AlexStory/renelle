// ast/ast.go

package ast

import (
	"bytes"
	"renelle/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
	Comments() []string
	AddComment(string)
}

type Expression interface {
	Node
	expressionNode()
	Comments() []string
	AddComment(string)
}

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token    token.Token // the token.LET token
	Name     *Identifier
	Value    Expression
	comments []string
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) Comments() []string   { return ls.comments }
func (ls *LetStatement) AddComment(c string)  { ls.comments = append(ls.comments, c) }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}

type Identifier struct {
	Token    token.Token // the token.IDENT token
	Value    string
	comments []string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) Comments() []string   { return i.comments }
func (i *Identifier) AddComment(c string)  { i.comments = append(i.comments, c) }

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
	comments    []string
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) Comments() []string   { return rs.comments }
func (rs *ReturnStatement) AddComment(c string)  { rs.comments = append(rs.comments, c) }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression

	comments []string
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) Comments() []string   { return es.comments }
func (es *ExpressionStatement) AddComment(c string)  { es.comments = append(es.comments, c) }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64

	comments []string
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) Comments() []string   { return il.comments }
func (il *IntegerLiteral) AddComment(c string)  { il.comments = append(il.comments, c) }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64

	comments []string
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }
func (fl *FloatLiteral) Comments() []string   { return fl.comments }
func (fl *FloatLiteral) AddComment(c string)  { fl.comments = append(fl.comments, c) }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression

	comments []string
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) Comments() []string   { return pe.comments }
func (pe *PrefixExpression) AddComment(c string)  { pe.comments = append(pe.comments, c) }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression

	comments []string
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) Comments() []string   { return ie.comments }
func (ie *InfixExpression) AddComment(c string)  { ie.comments = append(ie.comments, c) }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool

	comments []string
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) Comments() []string   { return b.comments }
func (b *Boolean) AddComment(c string)  { b.comments = append(b.comments, c) }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement

	comments []string
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) Comments() []string   { return ie.comments }
func (ie *IfExpression) AddComment(c string)  { ie.comments = append(ie.comments, c) }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement

	comments []string
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) Comments() []string   { return bs.comments }
func (bs *BlockStatement) AddComment(c string)  { bs.comments = append(bs.comments, c) }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement

	comments []string
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) Comments() []string   { return fl.comments }
func (fl *FunctionLiteral) AddComment(c string)  { fl.comments = append(fl.comments, c) }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("//")
	out.WriteString(strings.Join(params, " "))
	out.WriteString(" => ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type FunctionStatement struct {
	Token      token.Token // the 'fn' token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement

	comments []string
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) Comments() []string   { return fs.comments }
func (fs *FunctionStatement) AddComment(c string)  { fs.comments = append(fs.comments, c) }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral() + " ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, " "))
	out.WriteString(") ")
	out.WriteString(fs.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression

	comments []string
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) Comments() []string   { return ce.comments }
func (ce *CallExpression) AddComment(c string)  { ce.comments = append(ce.comments, c) }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

type AtomExpression struct {
	Token token.Token
	Value string

	comments []string
}

func (ae *AtomExpression) expressionNode()      {}
func (ae *AtomExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AtomExpression) Comments() []string   { return ae.comments }
func (ae *AtomExpression) AddComment(c string)  { ae.comments = append(ae.comments, c) }
func (ae *AtomExpression) String() string       { return ae.Value }
