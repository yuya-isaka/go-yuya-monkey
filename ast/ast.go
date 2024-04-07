package ast

import (
	"bytes"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Node interface {
	GetTokenContent() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

//--------------------

type Program struct {
	StatementArray []Statement
}

func (p *Program) GetTokenContent() string {
	if len(p.StatementArray) > 0 {
		return p.StatementArray[0].GetTokenContent()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.StatementArray {
		out.WriteString(s.String())
	}

	return out.String()
}

//--------------------

type LetStatement_1 struct {
	Token     token.Token // let
	IdentName *Identifier
	LetValue  Expression
}

func (ls *LetStatement_1) statementNode() {}
func (ls LetStatement_1) GetTokenContent() string {
	return ls.Token.Content
}
func (ls LetStatement_1) String() string {
	var out bytes.Buffer

	out.WriteString(ls.GetTokenContent() + " ")
	out.WriteString(ls.IdentName.String())
	out.WriteString(" = ")

	if ls.LetValue != nil {
		out.WriteString(ls.LetValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token      token.Token
	IdentValue string
}

func (i *Identifier) expressionNode() {}
func (i Identifier) GetTokenContent() string {
	return i.Token.Content
}
func (i Identifier) String() string {
	return i.IdentValue
}

type ReturnStatement_2 struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement_2) statementNode() {}
func (rs ReturnStatement_2) GetTokenContent() string {
	return rs.Token.Content
}
func (rs ReturnStatement_2) String() string {
	var out bytes.Buffer

	out.WriteString(rs.GetTokenContent() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement_3 struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement_3) statementNode() {}
func (es ExpressionStatement_3) GetTokenContent() string {
	return es.Token.Content
}
func (es ExpressionStatement_3) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Integer struct {
	Token        token.Token
	IntegerValue int64
}

func (i *Integer) expressionNode() {}
func (i Integer) GetTokenContent() string {
	return i.Token.Content
}
func (i Integer) String() string {
	return i.Token.Content
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) GetTokenContent() string {
	return p.Token.Content
}
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}
