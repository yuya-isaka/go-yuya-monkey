package ast

import (
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/token"
)

func TestString(t *testing.T) {
	input := "let myVar = anotherVar;"

	program := &ProgramNode{
		Statements: []Statement{
			&LetNode{
				Token: token.Token{Type: token.LET, Name: "let"},
				Name: &IdentNode{
					Token: token.Token{Type: token.IDENT, Name: "myVar"},
					Value: "myVar",
				},
				Value: &IdentNode{
					Token: token.Token{Type: token.IDENT, Name: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != input {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
