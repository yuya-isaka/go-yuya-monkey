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

func testLetStatementIs(t *testing.T, s ast.Statement, name string) bool {
	if s.GetTokenContent() != "let" {
		t.Errorf("s.TokenContent not 'let'. got=%q", s.GetTokenContent())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement_1)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.IdentName.IdentValue != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.IdentName.IdentValue)
		return false
	}

	if letStmt.IdentName.GetTokenContent() != name {
		t.Errorf("letStmt.Name.TokenContent() not '%s'. got=%s", name, letStmt.IdentName.GetTokenContent())
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
		returnStmt, ok := stmt.(*ast.ReturnStatement_2)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.GetTokenContent() != "return" {
			t.Errorf("returnStmt.TokenContent not 'return', got %q", returnStmt.GetTokenContent())
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.StatementArray) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.StatementArray))
	}

	stmt, ok := program.StatementArray[0].(*ast.ExpressionStatement_3)
	if !ok {
		t.Fatalf("program.StatementArray[0] is not ast.ExpressionStatement_3. got=%T", program.StatementArray[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.IdentValue != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.IdentValue)
	}
	if ident.GetTokenContent() != "foobar" {
		t.Errorf("ident.GetTokenContent not %s. got=%s", "foobar", ident.GetTokenContent())
	}
}

func TestIntegerContentExpression(t *testing.T) {
	input := "42;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// 正しい数の文が生成されたか
	if len(program.StatementArray) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.StatementArray))
	}

	// 最初の文がExpressionStatementだよね？
	stmt, ok := program.StatementArray[0].(*ast.ExpressionStatement_3)
	if !ok {
		t.Fatalf("program.StatementArray[0] is not ast.ExpressionStatement_3. got=%T", program.StatementArray[0])
	}

	// それの式はIntegerContentだよね？
	integerContent, ok := stmt.Expression.(*ast.IntegerContent)
	if !ok {
		t.Fatalf("exp not *ast.IntegerContent. got=%T", stmt.Expression)
	}
	// それの中身は42だよね？
	if integerContent.IntegerValue != 42 {
		t.Errorf("integerContent.IntegerValue not %d. got=%d", 42, integerContent.IntegerValue)
	}
	// トークンとして取得したら"42"だよね？
	if integerContent.GetTokenContent() != "42" {
		t.Errorf("integerContent.GetTokenContent not %s. got=%s", "42", integerContent.GetTokenContent())
	}

	// 上記のアサーションを設ける
}
