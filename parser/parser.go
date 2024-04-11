package parser

import (
	"fmt"
	"strconv"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/lexer"
	"github.com/yuya-isaka/go-yuya-monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	lex    *lexer.Lexer
	errors []string

	curT  token.Token
	peekT token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lex:    l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.ProgramNode {
	// defer untrace(trace("ParseProgram"))

	// 作る
	node := &ast.ProgramNode{
		Statements: []ast.Statement{},
	}

	// 文をトークンの最後まで
	for !p.curToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			msg := fmt.Sprintf("文がnilだぜ %T", stmt)
			p.errors = append(p.errors, msg)
		}
		node.Statements = append(node.Statements, stmt)

		// セミコロンで返ってきているはず
		// 式文なら文の末尾
		p.nextToken()
	}

	return node
}

func (p *Parser) parseStatement() ast.Statement {
	// defer untrace(trace("parseStatement"))

	switch p.curT.Type {
	case token.LET:
		return p.parseLet()

	case token.RETURN:
		return p.parseReturn()

	default:
		return p.parseES() // 式文
	}
}

func (p *Parser) parseLet() ast.Statement {
	// defer untrace(trace("parseLet"))

	node := &ast.LetNode{Token: p.curT}

	// let a = 3;
	//  ↑
	if !p.expectPeekToken(token.IDENT) {
		return nil
	}

	node.Name = &ast.IdentNode{Token: p.curT, Value: p.curT.Name}

	// let a = 3;
	//     ↑
	if !p.expectPeekToken(token.ASSIGN) {
		return nil
	}

	// let a = 3;
	//       ↑
	p.nextToken()

	// let a = 3;
	//         ↑
	node.Value = p.parseExpression(LOWEST)

	// どちらか
	// let a = 3;
	//         ↑
	// or
	// let a = 3;
	//          ↑
	if !p.curToken(token.SEMICOLON) {
		p.nextToken()
	}

	if !p.curToken(token.SEMICOLON) {
		msg := fmt.Sprintf("\";\" is nothing!!! token is %q", p.curT.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// セミコロンで返る
	return node
}

func (p *Parser) parseReturn() ast.Statement {
	// defer untrace(trace("parseReturn"))

	node := &ast.ReturnNode{Token: p.curT}

	p.nextToken()

	node.Value = p.parseExpression(LOWEST)

	if !p.curToken(token.SEMICOLON) {
		p.nextToken()
	}

	if !p.curToken(token.SEMICOLON) {
		msg := fmt.Sprintf("\";\" is nothing!!! token is %q", p.curT.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// セミコロンで返る
	return node
}

func (p *Parser) parseES() ast.Statement {
	// defer untrace(trace("parseES"))

	node := &ast.EsNode{Token: p.curT}
	node.Value = p.parseExpression(LOWEST)

	// 式文のセミコロンなしOK
	if p.peekToken(token.SEMICOLON) {
		p.nextToken()
	}

	// セミコロンか式文の最後で返る
	return node
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))

	var left ast.Expression

	switch p.curT.Type {
	case token.IDENT:
		left = p.parseIdent()

	case token.INT:
		left = p.parseInt()

	case token.BANG, token.MINUS:
		left = p.parsePrefix()

	case token.TRUE, token.FALSE:
		left = p.parseBool()

	case token.LPAREN:
		left = p.parseGroup()

	case token.IF:
		left = p.parseIf()

	case token.FUNCTION:
		left = p.parseFunction()

	default:
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curT.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	for !p.peekToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		switch p.peekT.Type {
		case token.EQ, token.NOT_EQ, token.LT, token.GT, token.PLUS, token.MINUS, token.ASTERISK, token.SLASH:
			p.nextToken()
			left = p.parseInfix(left)

		case token.LPAREN:
			p.nextToken()
			left = p.parseCall(left)

		default:
			msg := fmt.Sprintf("ここにくるのはなんだ %q", p.peekT.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return left
}

//--------------------

func (p *Parser) nextToken() {
	p.curT = p.peekT
	p.peekT = p.lex.NextToken()
}

func (p *Parser) curToken(t token.TokenType) bool {
	return p.curT.Type == t
}

func (p *Parser) peekToken(t token.TokenType) bool {
	return p.peekT.Type == t
}

func (p *Parser) expectPeekToken(t token.TokenType) bool {
	if p.peekToken(t) {
		p.nextToken()
		return true
	} else {
		msg := fmt.Sprintf("expected nexttoken to be %s, got %s instead", t, p.peekT.Type)
		p.errors = append(p.errors, msg)
		return false
	}
}

//--------------------

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekT.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curT.Type]; ok {
		return p
	}

	return LOWEST
}

//--------------------

func (p *Parser) Errors() []string {
	return p.errors
}

//--------------------

func (p *Parser) parseIdent() ast.Expression {
	// defer untrace(trace("parseIdent"))

	return &ast.IdentNode{Token: p.curT, Value: p.curT.Name}
}

func (p *Parser) parseInt() ast.Expression {
	// defer untrace(trace("parseInt"))

	value, err := strconv.ParseInt(p.curT.Name, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curT.Name)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntNode{Token: p.curT, Value: value}
}

func (p *Parser) parsePrefix() ast.Expression {
	// defer untrace(trace("parsePrefix"))

	node := &ast.PrefixNode{
		Token:    p.curT,
		Operator: p.curT.Name,
	}

	p.nextToken()

	node.Right = p.parseExpression(PREFIX)

	return node
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfix"))

	node := &ast.InfixNode{
		Token:    p.curT,
		Left:     left,
		Operator: p.curT.Name,
	}

	// ここが超大事
	// parseExpressionを呼び出すときの、precedenceをどうするかで大きく変わってくる。
	precedence := p.curPrecedence()

	p.nextToken()
	node.Right = p.parseExpression(precedence)

	return node
}

func (p *Parser) parseBool() ast.Expression {
	// defer untrace(trace("parseBool"))

	return &ast.BoolNode{Token: p.curT, Value: p.curToken(token.TRUE)}
}

func (p *Parser) parseGroup() ast.Expression {
	// defer untrace(trace("parseGroup"))

	p.nextToken()

	node := p.parseExpression(LOWEST)

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return node
}

func (p *Parser) parseIf() ast.Expression {
	// defer untrace(trace("parseIf"))

	node := &ast.IfNode{Token: p.curT}

	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}
	p.nextToken()

	node.Condition = p.parseExpression(LOWEST)

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	if !p.expectPeekToken(token.LBRACE) {
		return nil
	}

	node.Consequence = p.parseBlock() // 進む

	// elseを省略してもOK
	if p.peekToken(token.ELSE) {
		p.nextToken()

		if !p.expectPeekToken(token.LBRACE) {
			return nil
		}

		node.Alternative = p.parseBlock()
	}

	return node
}

// 返る先の型が指定されているから、返り値の型はast.Statementではなく*ast.BlockNode
func (p *Parser) parseBlock() *ast.BlockNode {
	// defer untrace(trace("parseBlock"))

	node := &ast.BlockNode{
		Token:      p.curT,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	// 文をトークンの最後まで（"}"を忘れた場合、次々に"}"が出てくるまで、トークンを勧めながら文を読もうとしてしまうのでそれを避けるために、token.EOFでも確認）
	for !p.curToken(token.RBRACE) && !p.curToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			msg := fmt.Sprintf("文がnilだぜ %T", stmt)
			p.errors = append(p.errors, msg)
		}

		node.Statements = append(node.Statements, stmt)

		// セミコロンで帰ってくるはず
		// 式文なら文の末尾かも
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFunction() ast.Expression {
	// defer untrace(trace("parseFunction"))

	node := &ast.FunctionNode{Token: p.curT}

	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}

	node.Parameters = p.parseParameters()

	if !p.expectPeekToken(token.LBRACE) {
		return nil
	}

	node.Body = p.parseBlock()

	return node
}

// 返る先の型が指定されているから、返り値の型はast.Expressionではなく*ast.IdentNode
func (p *Parser) parseParameters() []*ast.IdentNode {
	// defer untrace(trace("parseParameters"))

	nodes := []*ast.IdentNode{}

	if p.peekToken(token.RPAREN) {
		p.nextToken()
		return nodes
	}

	p.nextToken()

	node := &ast.IdentNode{Token: p.curT, Value: p.curT.Name}
	nodes = append(nodes, node)

	for p.peekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		node := &ast.IdentNode{Token: p.curT, Value: p.curT.Name}
		nodes = append(nodes, node)
	}

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return nodes
}

func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	// defer untrace(trace("parseCall"))

	node := &ast.CallNode{Token: p.curT, Function: function}
	node.Arguments = p.parseArguments()
	return node
}

func (p *Parser) parseArguments() []ast.Expression {
	// defer untrace(trace("parseArguments"))

	nodes := []ast.Expression{}

	if p.peekToken(token.RPAREN) {
		p.nextToken()
		return nodes
	}

	p.nextToken()
	nodes = append(nodes, p.parseExpression(LOWEST))

	for p.peekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		nodes = append(nodes, p.parseExpression(LOWEST))
	}

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return nodes
}
