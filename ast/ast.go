package ast

import (
	"bytes"
	"strings"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

// ---------------------------------

type ProgramNode struct {
	Statements []Statement
}

func (p ProgramNode) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

//--------------------

type LetNode struct {
	Token token.Token // let
	Name  *IdentNode
	Value Expression
}

func (l LetNode) statement() {}
func (l LetNode) String() string {
	var out bytes.Buffer

	out.WriteString(l.Token.Name + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ---------------------------------

type IdentNode struct {
	Token token.Token
	Value string
}

func (i IdentNode) expression() {}
func (i IdentNode) String() string {
	return i.Value
}

// ---------------------------------

type ReturnNode struct {
	Token token.Token
	Value Expression
}

func (rs ReturnNode) statement() {}
func (rs ReturnNode) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Token.Name + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ---------------------------------

// Expression Statement
type EsNode struct {
	Token token.Token
	Value Expression
}

func (es EsNode) statement() {}
func (es EsNode) String() string {
	if es.Value != nil {
		return es.Value.String()
	}
	return ""
}

// ---------------------------------

type IntNode struct {
	Token token.Token
	Value int64
}

func (i IntNode) expression() {}
func (i IntNode) String() string {
	return i.Token.Name
}

// ---------------------------------

type PrefixNode struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p PrefixNode) expression() {}
func (p PrefixNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// ---------------------------------

type InfixNode struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i InfixNode) expression() {}
func (i InfixNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

type BoolNode struct {
	Token token.Token
	Value bool
}

func (b BoolNode) expression()    {}
func (b BoolNode) String() string { return b.Token.Name }

// ------------------------

type IfNode struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockNode
	Alternative *BlockNode
}

func (ie IfNode) expression() {}
func (ie IfNode) String() string {
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

type BlockNode struct {
	Token      token.Token
	Statements []Statement
}

func (bs BlockNode) statement() {}
func (bs BlockNode) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionNode struct {
	Token      token.Token // 'fn'
	Parameters []*IdentNode
	Body       *BlockNode
}

func (f FunctionNode) expression() {}
func (f FunctionNode) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.Token.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

type CallNode struct {
	Token     token.Token
	Function  Expression // Identifier or Function
	Arguments []Expression
}

func (c CallNode) expression() {}
func (c CallNode) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringNode struct {
	Token token.Token
	Value string
}

func (s StringNode) expression()    {}
func (s StringNode) String() string { return s.Token.Name }
