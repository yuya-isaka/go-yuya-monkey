package ast

import (
	"bytes"
	"strings"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Node interface {
	String() string // ASTの表現を見える化！
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
	Token token.Token // 先頭のトークン
	Name  *IdentNode  // 変数名
	Value Expression  // 中の式
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

type ReturnNode struct {
	Token token.Token // 先頭のトークン
	Value Expression  // 返す式
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

// Expression Statement
type EsNode struct {
	Token token.Token // 先頭のトークン
	Value Expression  // 持つ式
}

func (e EsNode) statement() {}
func (e EsNode) String() string {

	if e.Value != nil {
		return e.Value.String()
	}
	return ""
}

type BlockNode struct {
	Token      token.Token
	Statements []Statement
}

func (b BlockNode) statement() {}
func (b BlockNode) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ---------------------------------

// 変数名
type IdentNode struct {
	Token token.Token // 先頭のトークン
	Value string
}

func (i IdentNode) expression() {}
func (i IdentNode) String() string {
	return i.Value
}

type IntNode struct {
	Token token.Token // 先頭のトークン
	Value int64       // 持つ値
}

func (i IntNode) expression() {}
func (i IntNode) String() string {
	return i.Token.Name
}

type PrefixNode struct {
	Token    token.Token // 先頭のトークン
	Operator string      // オペレータ
	Right    Expression  // 右辺の式
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

type InfixNode struct {
	Token    token.Token // 先頭のトークン
	Left     Expression  // 左辺
	Operator string      // オペレータ
	Right    Expression  // 右辺
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
	Token token.Token // 先頭のトークン
	Value bool        // 値
}

func (b BoolNode) expression()    {}
func (b BoolNode) String() string { return b.Token.Name }

// ------------------------

type IfNode struct {
	Token       token.Token // 先頭のトークン
	Condition   Expression  // 条件式
	Consequence *BlockNode  // ブロックノード
	Alternative *BlockNode  // ブロックノード
}

func (i IfNode) expression() {}
func (i IfNode) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

type FunctionNode struct {
	Token      token.Token  // 'fn'トークン、先頭のトークン
	Parameters []*IdentNode // 変数の配列
	Body       *BlockNode   // ブロックノード
}

func (f FunctionNode) expression() {}
func (f FunctionNode) String() string {
	var out bytes.Buffer

	params := make([]string, 0, len(f.Parameters))
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
	Token     token.Token  // 先頭のトークン
	Function  Expression   // Identifier or Function
	Arguments []Expression // 式の配列（先頭から評価）
}

func (c CallNode) expression() {}
func (c CallNode) String() string {
	var out bytes.Buffer

	args := make([]string, 0, len(c.Arguments))
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
	Token token.Token // 先頭のトークン
	Value string      // 値
}

func (s StringNode) expression()    {}
func (s StringNode) String() string { return s.Token.Name }
