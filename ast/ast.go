package ast

import (
	"bytes"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Node interface {
	GetTokenLiteral() string
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

func (p *Program) GetTokenLiteral() string {
	if len(p.StatementArray) > 0 {
		return p.StatementArray[0].GetTokenLiteral()
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
func (ls LetStatement_1) GetTokenLiteral() string {
	return ls.Token.Literal
}
func (ls LetStatement_1) String() string {
	var out bytes.Buffer

	out.WriteString(ls.GetTokenLiteral() + " ")
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
func (i Identifier) GetTokenLiteral() string {
	return i.Token.Literal
}
func (i Identifier) String() string {
	return i.IdentValue
}

type ReturnStatement_2 struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement_2) statementNode() {}
func (rs ReturnStatement_2) GetTokenLiteral() string {
	return rs.Token.Literal
}
func (rs ReturnStatement_2) String() string {
	var out bytes.Buffer

	out.WriteString(rs.GetTokenLiteral() + " ")

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
func (es ExpressionStatement_3) GetTokenLiteral() string {
	return es.Token.Literal
}
func (es ExpressionStatement_3) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
