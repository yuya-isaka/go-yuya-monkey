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

type LetStatement struct {
	Token token.Token // let
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()         {}
func (ls LetStatement) GetTokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()        {}
func (i Identifier) GetTokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()          {}
func (rs *ReturnStatement) GetTokenLiteral() string { return rs.Token.Literal }
