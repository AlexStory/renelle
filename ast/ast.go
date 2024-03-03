// ast/ast.go

package ast

import (
	"bytes"
	"renelle/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	T() token.Token
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

func (p Program) T() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].T()
	} else {
		return token.Token{}
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
	Left     Expression
	Value    Expression
	comments []string
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) T() token.Token       { return ls.Token }
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) Comments() []string   { return ls.comments }
func (ls *LetStatement) AddComment(c string)  { ls.comments = append(ls.comments, c) }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Left.String())
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
func (i *Identifier) T() token.Token       { return i.Token }
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
func (rs *ReturnStatement) T() token.Token       { return rs.Token }
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
func (es *ExpressionStatement) T() token.Token       { return es.Token }
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
func (il *IntegerLiteral) T() token.Token       { return il.Token }
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
func (fl *FloatLiteral) T() token.Token       { return fl.Token }
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }
func (fl *FloatLiteral) Comments() []string   { return fl.comments }
func (fl *FloatLiteral) AddComment(c string)  { fl.comments = append(fl.comments, c) }

type StringLiteral struct {
	Token    token.Token
	Value    string
	comments []string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) T() token.Token       { return sl.Token }
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) Comments() []string   { return sl.comments }
func (sl *StringLiteral) AddComment(c string)  { sl.comments = append(sl.comments, c) }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
	comments []string
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) T() token.Token       { return al.Token }
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) Comments() []string   { return al.comments }
func (al *ArrayLiteral) AddComment(c string)  { al.comments = append(al.comments, c) }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, " "))
	out.WriteString("]")
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression

	comments []string
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) T() token.Token       { return pe.Token }
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
func (ie *InfixExpression) T() token.Token       { return ie.Token }
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
func (b *Boolean) T() token.Token       { return b.Token }
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
func (ie *IfExpression) T() token.Token       { return ie.Token }
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
func (bs *BlockStatement) T() token.Token       { return bs.Token }
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
func (fl *FunctionLiteral) T() token.Token       { return fl.Token }
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
func (fs *FunctionStatement) T() token.Token       { return fs.Token }
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
func (ce *CallExpression) T() token.Token       { return ce.Token }
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

type AtomLiteral struct {
	Token token.Token
	Value string

	comments []string
}

func (ae *AtomLiteral) expressionNode()      {}
func (ae *AtomLiteral) T() token.Token       { return ae.Token }
func (ae *AtomLiteral) TokenLiteral() string { return ae.Token.Literal }
func (ae *AtomLiteral) Comments() []string   { return ae.comments }
func (ae *AtomLiteral) AddComment(c string)  { ae.comments = append(ae.comments, c) }
func (ae *AtomLiteral) String() string       { return ae.Value }

type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression

	comments []string
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) T() token.Token       { return ie.Token }
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) Comments() []string   { return ie.comments }
func (ie *IndexExpression) AddComment(c string)  { ie.comments = append(ie.comments, c) }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" @ ")
	out.WriteString(ie.Index.String())
	out.WriteString(")")
	return out.String()
}

type PropertyAccessExpression struct {
	Token token.Token // The '.' token
	Left  Expression
	Right *Identifier

	comments []string
}

func (pae *PropertyAccessExpression) expressionNode()      {}
func (pae *PropertyAccessExpression) T() token.Token       { return pae.Token }
func (pae *PropertyAccessExpression) TokenLiteral() string { return pae.Token.Literal }
func (pae *PropertyAccessExpression) Comments() []string   { return pae.comments }
func (pae *PropertyAccessExpression) AddComment(c string)  { pae.comments = append(pae.comments, c) }
func (pae *PropertyAccessExpression) String() string {
	return pae.Left.String() + "." + pae.Right.String()
}

type ApplyExpression struct {
	Token token.Token // The '$' token
	Left  Expression
	Right Expression

	comments []string
}

func (ae *ApplyExpression) expressionNode()      {}
func (ae *ApplyExpression) T() token.Token       { return ae.Token }
func (ae *ApplyExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *ApplyExpression) String() string       { return ae.Left.String() + " $ " + ae.Right.String() }
func (ae *ApplyExpression) Comments() []string   { return ae.comments }
func (ae *ApplyExpression) AddComment(c string)  { ae.comments = append(ae.comments, c) }

type TupleLiteral struct {
	Token    token.Token // The '(' token
	Elements []Expression

	comments []string
}

func (te *TupleLiteral) expressionNode()      {}
func (te *TupleLiteral) T() token.Token       { return te.Token }
func (te *TupleLiteral) TokenLiteral() string { return te.Token.Literal }
func (te *TupleLiteral) Comments() []string   { return te.comments }
func (te *TupleLiteral) AddComment(c string)  { te.comments = append(te.comments, c) }
func (te *TupleLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range te.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(elements, " "))
	out.WriteString(")")
	return out.String()
}

type MapLiteral struct {
	Token token.Token // The '{' token
	Pairs map[Expression]Expression

	comments []string
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) T() token.Token       { return ml.Token }
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MapLiteral) Comments() []string   { return ml.comments }
func (ml *MapLiteral) AddComment(c string)  { ml.comments = append(ml.comments, c) }
func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range ml.Pairs {
		pairs = append(pairs, key.String()+" = "+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, " "))
	out.WriteString("}")
	return out.String()
}

type CondExpression struct {
	Token        token.Token // The 'cond' token
	Conditions   []Expression
	Consequences []*BlockStatement

	comments []string
}

func (ce *CondExpression) expressionNode()      {}
func (ce *CondExpression) T() token.Token       { return ce.Token }
func (ce *CondExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CondExpression) Comments() []string   { return ce.comments }
func (ce *CondExpression) AddComment(c string)  { ce.comments = append(ce.comments, c) }
func (ce *CondExpression) String() string {
	var out bytes.Buffer
	out.WriteString("cond")
	for i, cond := range ce.Conditions {
		out.WriteString(cond.String())
		out.WriteString(ce.Consequences[i].String())
	}
	return out.String()
}
