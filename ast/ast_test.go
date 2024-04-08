package ast

import (
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		StatementArray: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Content: "let"},
				LetName: &Identifier{
					Token:      token.Token{Type: token.IDENT, Content: "myVar"},
					IdentValue: "myVar",
				},
				LetExpression: &Identifier{
					Token:      token.Token{Type: token.IDENT, Content: "anotherVar"},
					IdentValue: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
