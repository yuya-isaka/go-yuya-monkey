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
)

type TokenType string

type Token struct {
	Type    TokenType
	Content string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookKeywordToken(content string) TokenType {
	if tok, ok := keywords[content]; ok {
		return tok
	}
	return IDENT
}
