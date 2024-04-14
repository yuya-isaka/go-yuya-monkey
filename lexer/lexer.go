package lexer

import (
	"github.com/yuya-isaka/go-yuya-monkey/token"
)

type Lexer struct {
	input string // ソースコード全部
	pos   int    // 読んでいる場所
	ch    byte   // 読んでいる場所のバイト
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

	// 空白文字, 改行, タブ, キャリッジリターンを無視
	for l.ch == '\n' || l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.nextPos()
	}

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
	case 0: // EOF==0, 整数0==48
		tok = newToken(token.EOF, string(l.ch))
	case '"':
		pos := l.pos + 1
		for {
			l.nextPos()
			// EOFで判断しないと"出るまで永遠に終わらない
			// 「何かが出たら終わる」っていう条件分岐をするときは、対象のものが出ない時のことを考える
			if l.ch == '"' || l.ch == 0 {
				break
			}
		}
		tok = newToken(token.STRING, l.input[pos:l.pos])
	default:
		switch {
		case isLetter(l.ch):
			pos := l.pos
			// ーーじゃなかったら終わり系ははっきりしている
			for isLetter(l.ch) {
				l.nextPos()
			}
			// 次の文字まで進んでしまっているからここでリターン
			// 文字だったら「キーワード」チェック。キーワードか変数かここじゃわからん
			return newToken(token.LookKeyword(l.input[pos:l.pos]), l.input[pos:l.pos])
		case isNumber(l.ch):
			pos := l.pos
			// ーーじゃなかったら終わり系ははっきりしている
			for isNumber(l.ch) {
				l.nextPos()
			}
			// 次の文字まで進んでしまっているからここでリターン
			return newToken(token.INT, l.input[pos:l.pos])
		default:
			// おかしい
			tok = newToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.nextPos()

	return tok
}

// -----------------------------------------------------------------

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

func (l Lexer) peek() byte {
	peekPos := l.pos + 1
	// 先を見るときは境界チェックを気をつけて
	if peekPos >= len(l.input) {
		return 0
	}
	return l.input[peekPos]
}
