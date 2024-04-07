package ast

import (
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		StatementArray: []Statement{
			&LetStatement_1{
				Token: token.Token{Type: token.LET, Literal: "let"},
				IdentName: &Identifier{
					Token:      token.Token{Type: token.IDENT, Literal: "myVar"},
					IdentValue: "myVar",
				},
				LetValue: &Identifier{
					Token:      token.Token{Type: token.IDENT, Literal: "anotherVar"},
					IdentValue: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
