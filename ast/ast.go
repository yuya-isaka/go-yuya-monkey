package ast

import (
	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Node interface {
	GetTokenLiteral() string
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

type LetStatement_1 struct {
	Token token.Token // let
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement_1) statementNode()         {}
func (ls LetStatement_1) GetTokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()        {}
func (i Identifier) GetTokenLiteral() string { return i.Token.Literal }

type ReturnStatement_2 struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement_2) statementNode()         {}
func (rs ReturnStatement_2) GetTokenLiteral() string { return rs.Token.Literal }

type ExpressionStatement_3 struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement_3) expressionNode() {}
func (es ExpressionStatement_3) GetTokenLiteral() string {
	return es.Token.Literal
}
