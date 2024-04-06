package lexer

import (
	"github.com/yiuya-isaka/go-yuya-monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.eatChar()
	return l
}

func (l *Lexer) eatChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var newTok token.Token

	l.eatWhitespace()

	switch l.ch {
	case '=':
		newTok = newToken(token.ASSIGN, l.ch)
	case '+':
		newTok = newToken(token.PLUS, l.ch)
	case ',':
		newTok = newToken(token.COMMA, l.ch)
	case ';':
		newTok = newToken(token.SEMICOLON, l.ch)
	case '(':
		newTok = newToken(token.LPAREN, l.ch)
	case ')':
		newTok = newToken(token.RPAREN, l.ch)
	case '{':
		newTok = newToken(token.LBRACE, l.ch)
	case '}':
		newTok = newToken(token.RBRACE, l.ch)
	case 0:
		newTok = newToken(token.EOF, l.ch)
	default:
		if isLetter(l.ch) {
			newTok.Literal = l.readIdentifier()
			newTok.Type = token.LookKeywordToken(newTok.Literal)
			return newTok
		}
		newTok = newToken(token.ILLEGAL, l.ch)
	}

	l.eatChar()

	return newTok
}

func newToken(tt token.TokenType, literal byte) token.Token {
	return token.Token{Type: tt, Literal: string(literal)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.eatChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) eatWhitespace() {
	for l.ch == '\n' || l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.eatChar()
	}
}
