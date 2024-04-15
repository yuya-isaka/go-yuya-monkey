package parser

// 構文解析とは、Token構造体をNodeインタフェースを満たした構造体にすること

import (
	"fmt"
	"strconv"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/lexer"
	"github.com/yuya-isaka/go-yuya-monkey/token"
)

// 優先順位
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

// 優先順位表
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
	token.LBRACKET: INDEX,
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

// --------------------------------------------------------------------------

func (p *Parser) ParseProgram() *ast.ProgramNode {

	node := &ast.ProgramNode{
		Statements: []ast.Statement{},
	}

	// 文をトークンの最後まで
	for !p.curToken(token.EOF) {
		stmt := p.parseStatement()

		// 構文解析のエラーは、エラー配列にためてnilで返ってくるようにしている
		if stmt == nil {
			msg := fmt.Sprintf("文がnilだぜ %T", stmt)
			p.errors = append(p.errors, msg)
		} else {
			node.Statements = append(node.Statements, stmt)
		}

		// [セミコロン] or [式文なら文の末尾]で返ってきているはず
		p.nextToken()
	}

	return node
}

// 返る先の型が指定されているから、返り値の型はast.Statementではなく*ast.BlockNode
func (p *Parser) parseBlock() *ast.BlockNode {

	node := &ast.BlockNode{
		Token:      p.curT,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	// 文をトークンの最後まで（"}"を忘れた場合、次々に"}"が出てくるまで、トークンを勧めながら文を読もうとしてしまうのでそれを避けるために、token.EOFでも確認）
	for !p.curToken(token.RBRACE) && !p.curToken(token.EOF) {

		stmt := p.parseStatement()

		// 構文解析のエラーは、エラー配列にためてnilで返ってくるようにしている
		if stmt == nil {
			msg := fmt.Sprintf("文がnilだぜ %T", stmt)
			p.errors = append(p.errors, msg)
		} else {
			node.Statements = append(node.Statements, stmt)
		}

		// [セミコロン] or [式文なら文の末尾]で返ってきているはず
		p.nextToken()
	}

	return node
}

// --------------------------------------------------------------------------

func (p *Parser) parseStatement() ast.Statement {

	// トークンで判断するPratt構文解析
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
	// セミコロンのわけがない
	if p.curToken(token.SEMICOLON) {
		p.errors = append(p.errors, "\";\" is wrong!!!")
		return nil
	}

	p.nextToken()

	// letは";"が必須です
	if !p.curToken(token.SEMICOLON) {
		msg := fmt.Sprintf("\";\" is nothing!!! token is %q", p.curT.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// セミコロンで返る
	return node
}

func (p *Parser) parseReturn() ast.Statement {
	node := &ast.ReturnNode{Token: p.curT}

	p.nextToken()

	node.Value = p.parseExpression(LOWEST)

	// セミコロンのわけがない
	if p.curToken(token.SEMICOLON) {
		p.errors = append(p.errors, "\";\" is wrong!!!")
		return nil
	}

	p.nextToken()

	// returnは";"が必須です
	if !p.curToken(token.SEMICOLON) {
		msg := fmt.Sprintf("\";\" is nothing!!! token is %q", p.curT.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// セミコロンで返る
	return node
}

func (p *Parser) parseES() ast.Statement {
	node := &ast.EsNode{Token: p.curT}
	node.Value = p.parseExpression(LOWEST)

	// セミコロンのわけがない
	if p.curToken(token.SEMICOLON) {
		p.errors = append(p.errors, "\";\" is wrong!!!")
		return nil
	}

	// 式文のセミコロンなしOK
	if p.peekToken(token.SEMICOLON) {
		p.nextToken()
	}

	// [セミコロン] or [式文なら文の末尾]
	return node
}

// --------------------------------------------------------------------------

// precedenceに入っている優先順位は、真左の優先順位
// parseExpressionは、セミコロンの手前で必ず終わる（セミコロンあるないに関わらず最後で終わる）
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var left ast.Expression

	// トークンで呼び出す関数を決める（Pratt構文解析）
	// 下の関数たちはギリギリまで進める
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

	case token.STRING:
		left = p.parseString()

	case token.LBRACKET:
		left = p.parseArray()

	default:
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curT.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// 左辺 vs 右辺
	for !p.peekToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		switch p.peekT.Type {
		case token.EQ, token.NOT_EQ, token.LT, token.GT, token.PLUS, token.MINUS, token.ASTERISK, token.SLASH:
			p.nextToken()
			left = p.parseInfix(left)

		// 関数呼び出し
		case token.LPAREN:
			p.nextToken()
			left = p.parseCall(left)

		// 関数の後の[は、関数の評価された後のleftが入ってくる。それを配列の左辺として使う
		// 基本関数の左辺は変数しかこなくて、その場合precedenceは変数の優先順位(つまりLOWEST)
		// 関数より優先順位が低くても問題にはならない
		// 他のやつは問題になる（問題というかこれはどういう設計にするかという話でもある）
		// 			*が左にある場合（3 * [1,2,3][1]）は、*より優先順位が低いと左に引っ張られる
		// 			-や!が左にある場合（-[2,3,4][2]）は、-より優先順位が低いと左に引っ張られる
		case token.LBRACKET:
			p.nextToken()
			left = p.parseIndex(left)

		default:
			msg := fmt.Sprintf("Infix Error: %T (%+v)", p.peekT, p.peekT)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return left
}

// --------------------------------------------------------------------------

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

	// 表にないやつ一番下
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curT.Type]; ok {
		return p
	}

	// 表にないやつ一番下
	return LOWEST
}

//--------------------

func (p *Parser) Errors() []string {
	return p.errors
}

//--------------------

func (p *Parser) parseIdent() ast.Expression {
	return &ast.IdentNode{Token: p.curT, Value: p.curT.Name}
}

func (p *Parser) parseInt() ast.Expression {
	value, err := strconv.ParseInt(p.curT.Name, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curT.Name)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntNode{Token: p.curT, Value: value}
}

func (p *Parser) parsePrefix() ast.Expression {

	node := &ast.PrefixNode{
		Token:    p.curT,
		Operator: p.curT.Name,
	}

	p.nextToken()

	// 左がPREFIXだよーって投げる
	node.Right = p.parseExpression(PREFIX)

	return node
}

// これはトークンが真ん中の状態で呼ばれる
func (p *Parser) parseInfix(left ast.Expression) ast.Expression {

	node := &ast.InfixNode{
		Token:    p.curT,
		Left:     left,
		Operator: p.curT.Name,
	}

	// ここが超大事
	// parseExpressionを呼び出すときの、precedenceをどうするかで大きく変わってくる。
	precedence := p.curPrecedence()
	p.nextToken()
	// ここで式を再度呼ぶの痺れるな〜
	// ここが痺れるポイント覚えとけよー
	// ここでこのInfixのprecedence vs 次のpeekPrecedenceになってくる
	// 例えば 3 * [1,2,3,4][2]の場合、precedenceは*の優先順位 nextTokenは[で渡される
	// で、[1,2,3,4]が評価された後、forループに入って[が評価される
	//		ここで[の優先順位が低かったら、*に負けて、この右辺（[1,2,3,4]）は*に吸い込まれる
	//		でも[の優先順位が高かったら、[が勝って、*の右辺になるはずだったものは[の方に吸い込まれる（さらに右に）
	node.Right = p.parseExpression(precedence)

	return node
}

func (p *Parser) parseBool() ast.Expression {

	return &ast.BoolNode{Token: p.curT, Value: p.curToken(token.TRUE)}

}

func (p *Parser) parseGroup() ast.Expression {

	p.nextToken()

	node := p.parseExpression(LOWEST)

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return node
}

func (p *Parser) parseIf() ast.Expression {

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

func (p *Parser) parseFunction() ast.Expression {

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

	node := &ast.CallNode{Token: p.curT, Function: function}
	node.Arguments = p.parseExpressions(token.RPAREN)
	return node
}

func (p *Parser) parseArray() ast.Expression {
	array := &ast.ArrayNode{Token: p.curT}
	array.Values = p.parseExpressions(token.RBRACKET)
	return array
}

// callとarray で使っているよ
// 一般化されている
func (p *Parser) parseExpressions(end token.TokenType) []ast.Expression {

	nodes := []ast.Expression{}

	// 要素なしでreturn
	// endチェック
	if p.peekToken(end) {
		// ギリギリまで進めるよ
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

	// nilでreturn
	// endチェック
	if !p.expectPeekToken(end) {
		return nil
	}

	return nodes
}

func (p *Parser) parseString() ast.Expression {
	return &ast.StringNode{Token: p.curT, Value: p.curT.Name}
}

func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	// 呼ばれるときは先頭
	node := &ast.IndexNode{Token: p.curT, Left: left}
	p.nextToken()
	node.Index = p.parseExpression(LOWEST)
	// parseExpressionの後は、その式の最後で終わっている

	// 間違っている場合
	if !p.expectPeekToken(token.RBRACKET) {
		return nil
	}

	return node
}
