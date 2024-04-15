package token

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "="
	PLUS      = "+"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	BANG      = "!"
	MINUS     = "-"
	SLASH     = "/"
	ASTERISK  = "*"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	STRING    = "STRING"
	LBRACKET  = "["
	RBRACKET  = "]"
)

type TokenType string

type Token struct {
	Type TokenType // 型
	Name string    // 名前
}

// キーワードたち
// 名前で型を返す(名前→型)
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookKeyword(name string) TokenType {
	if tok, ok := keywords[name]; ok {
		return tok
	}
	// キーワードじゃなかったら変数
	return IDENT
}
