package parser

import (
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.StatementArray) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.StatementArray))
	}

	tests := []struct {
		expectIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.StatementArray[i]
		if !testLetStatementIs(t, stmt, tt.expectIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatementIs(t *testing.T, s ast.Statement, name string) bool {
	if s.GetTokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.GetTokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.GetTokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.GetTokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.StatementArray) != 3 {
		t.Fatalf("program.Statements does not contain 3 statemetns. got=%d", len(program.StatementArray))
	}

	for _, stmt := range program.StatementArray {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.GetTokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.GetTokenLiteral())
		}
	}
}
