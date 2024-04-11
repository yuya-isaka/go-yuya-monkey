package lexer

import (
	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Lexer struct {
	input string
	pos   int
	ch    byte
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		ch:    input[0],
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.ignoreSpace()

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			ch := l.ch
			l.nextPos()
			tok = newToken(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.ASSIGN, string(l.ch))
		}
	case '+':
		tok = newToken(token.PLUS, string(l.ch))
	case ',':
		tok = newToken(token.COMMA, string(l.ch))
	case ';':
		tok = newToken(token.SEMICOLON, string(l.ch))
	case '(':
		tok = newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = newToken(token.RPAREN, string(l.ch))
	case '{':
		tok = newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = newToken(token.RBRACE, string(l.ch))
	case '!':
		if l.peek() == '=' {
			ch := l.ch
			l.nextPos()
			tok = newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.BANG, string(l.ch))
		}
	case '-':
		tok = newToken(token.MINUS, string(l.ch))
	case '/':
		tok = newToken(token.SLASH, string(l.ch))
	case '*':
		tok = newToken(token.ASTERISK, string(l.ch))
	case '<':
		tok = newToken(token.LT, string(l.ch))
	case '>':
		tok = newToken(token.GT, string(l.ch))
	case 0:
		tok = newToken(token.EOF, string(l.ch))
	default:
		if isLetter(l.ch) {
			tok.Name = l.readIdent()
			tok.Type = token.LookKeyword(tok.Name)
			return tok
		} else if isNumber(l.ch) {
			tok.Name = l.readNumber()
			tok.Type = token.INT
			return tok
		}

		// おかしい
		tok = newToken(token.ILLEGAL, string(l.ch))
	}

	l.nextPos()

	return tok
}

func newToken(tt token.TokenType, name string) token.Token {
	return token.Token{Type: tt, Name: name}
}

func (l *Lexer) nextPos() {
	peekPos := l.pos + 1
	if peekPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[peekPos]
	}

	l.pos += 1
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.nextPos()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isNumber(l.ch) {
		l.nextPos()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) ignoreSpace() {
	for l.ch == '\n' || l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.nextPos()
	}
}

func (l Lexer) peek() byte {
	peekPos := l.pos + 1
	if peekPos >= len(l.input) {
		return 0
	}
	return l.input[peekPos]
}
