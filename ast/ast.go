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

type LetStatement struct {
	Token     token.Token // let
	IdentName *Identifier
	LetValue  Expression
}

func (ls *LetStatement) statementNode() {}
func (ls LetStatement) GetTokenContent() string {
	return ls.Token.Content
}
func (ls LetStatement) String() string {
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

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs ReturnStatement) GetTokenContent() string {
	return rs.Token.Content
}
func (rs ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.GetTokenContent() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es ExpressionStatement) GetTokenContent() string {
	return es.Token.Content
}
func (es ExpressionStatement) String() string {
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
func (p PrefixExpression) GetTokenContent() string {
	return p.Token.Content
}
func (p PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}
func (i InfixExpression) GetTokenContent() string {
	return i.Token.Content
}
func (i InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}
