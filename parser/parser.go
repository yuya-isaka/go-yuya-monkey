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
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFnMap map[token.TokenType]prefixParseFn
	infixParseFnMap  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFnMap = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerContent)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)  // 進める
	p.registerPrefix(token.MINUS, p.parsePrefixExpression) // 進める

	p.infixParseFnMap = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)     // 進める
	p.registerInfix(token.MINUS, p.parseInfixExpression)    // 進める
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) // 進める
	p.registerInfix(token.SLASH, p.parseInfixExpression)    // 進める
	p.registerInfix(token.EQ, p.parseInfixExpression)       // 進める
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)   // 進める
	p.registerInfix(token.LT, p.parseInfixExpression)       // 進める
	p.registerInfix(token.GT, p.parseInfixExpression)       // 進める

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	defer untrace(trace("ParseProgram"))

	program := &ast.Program{}
	program.StatementArray = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.StatementArray = append(program.StatementArray, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	defer untrace(trace("parseStatement"))

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	defer untrace(trace("parseLetStatement"))

	letstmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectTokenIs(token.IDENT) {
		return nil
	}

	letstmt.IdentName = &ast.Identifier{Token: p.curToken, IdentValue: p.curToken.Content}

	if !p.expectTokenIs(token.ASSIGN) {
		return nil
	}

	// TODO
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letstmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	defer untrace(trace("parseReturnStatement"))

	returnstmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnstmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))

	esstmt := &ast.ExpressionStatement{Token: p.curToken}
	esstmt.Expression = p.parseExpression(LOWEST)
	// 式文のセミコロンなしはエラーにしない。
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return esstmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFnMap[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// prefixは左結合がめちゃくちゃ強いのか (右にあるものの方が優先度が高い)
	// ってのを最初にprefixを処理するという展開にすることで、自然に表現している
	// 右にあるものの方が優先度が高いっていう、位置による優先度は、こんな感じで最初に処理するようにすれば、再帰的に必ず右にあるやつから処理されるのか。。。面白いな
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFnMap[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		// 1 +    2 +
		p.nextToken()
		// + 2    + 3
		leftExp = infix(leftExp)
		// 2 +    3 ;
	}

	return leftExp
}

//--------------------

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectTokenIs(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

//--------------------

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

//--------------------

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected nexttoken to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

//--------------------

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFnMap[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFnMap[tt] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	defer untrace(trace("parseIdentifier"))

	return &ast.Identifier{Token: p.curToken, IdentValue: p.curToken.Content}
}

func (p *Parser) parseIntegerContent() ast.Expression {
	defer untrace(trace("parseIntegerContent"))

	integerValue, err := strconv.ParseInt(p.curToken.Content, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Content)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.Integer{Token: p.curToken, IntegerValue: integerValue}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))

	prefixExpression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Content,
	}

	p.nextToken()

	prefixExpression.Right = p.parseExpression(PREFIX)

	return prefixExpression
}

func (p *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))

	infixExpression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     leftExpression,
		Operator: p.curToken.Content,
	}

	// ここが超大事
	// parseExpressionを呼び出すときの、precedenceをどうするかで大きく変わってくる。
	precedence := p.curPrecedence()

	// + 2    + 3
	p.nextToken()
	// 2 +    3 ;
	infixExpression.Right = p.parseExpression(precedence)

	return infixExpression
}

// 1 + 3 + 4
// LOWEST 1 + 3
//
