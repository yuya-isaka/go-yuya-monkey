package lexer

import (
	"github.com/yuya-isaka/go-yuya-monkey/token"
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
		if l.peek() == '=' {
			ch := l.ch
			l.eatChar()
			content := string(ch) + string(l.ch)
			newTok = newToken(token.EQ, content)
		} else {
			newTok = newToken(token.ASSIGN, string(l.ch))
		}
	case '+':
		newTok = newToken(token.PLUS, string(l.ch))
	case ',':
		newTok = newToken(token.COMMA, string(l.ch))
	case ';':
		newTok = newToken(token.SEMICOLON, string(l.ch))
	case '(':
		newTok = newToken(token.LPAREN, string(l.ch))
	case ')':
		newTok = newToken(token.RPAREN, string(l.ch))
	case '{':
		newTok = newToken(token.LBRACE, string(l.ch))
	case '}':
		newTok = newToken(token.RBRACE, string(l.ch))
	case '!':
		if l.peek() == '=' {
			ch := l.ch
			l.eatChar()
			content := string(ch) + string(l.ch)
			newTok = newToken(token.NOT_EQ, content)
		} else {
			newTok = newToken(token.BANG, string(l.ch))
		}
	case '-':
		newTok = newToken(token.MINUS, string(l.ch))
	case '/':
		newTok = newToken(token.SLASH, string(l.ch))
	case '*':
		newTok = newToken(token.ASTERISK, string(l.ch))
	case '<':
		newTok = newToken(token.LT, string(l.ch))
	case '>':
		newTok = newToken(token.GT, string(l.ch))
	case 0:
		newTok = newToken(token.EOF, string(l.ch))
	default:
		if isLetter(l.ch) {
			newTok.Content = l.readIdentifier()
			newTok.Type = token.LookKeywordToken(newTok.Content)
			return newTok
		} else if isNumber(l.ch) {
			newTok.Content = l.readNumber()
			newTok.Type = token.INT
			return newTok
		}
		newTok = newToken(token.ILLEGAL, string(l.ch))
	}

	l.eatChar()

	return newTok
}

func newToken(tt token.TokenType, content string) token.Token {
	return token.Token{Type: tt, Content: content}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.eatChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isNumber(l.ch) {
		l.eatChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) eatWhitespace() {
	for l.ch == '\n' || l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.eatChar()
	}
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
