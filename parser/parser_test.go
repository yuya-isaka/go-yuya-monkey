package parser

import (
	"fmt"
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/lexer"
	"github.com/yuya-isaka/go-yuya-monkey/token"
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

	letStmt, ok := s.(*ast.LetStatement)
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

//----------------------------------

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
		if returnStmt.GetTokenContent() != "return" {
			t.Errorf("returnStmt.TokenContent not 'return', got %q", returnStmt.GetTokenContent())
		}
	}
}

//----------------------------------

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

func testIntegerContent(t *testing.T, expression ast.Expression, integerValue int64) bool {
	integer, ok := expression.(*ast.Integer)
	if !ok {
		t.Errorf("expression no *ast.Integer. got=%T", expression)
		return false
	}

	if integer.IntegerValue != integerValue {
		t.Errorf("integer.IntegerValue not %d, got=%d", integerValue, integer.IntegerValue)
		return false
	}

	if integer.GetTokenContent() != fmt.Sprintf("%d", integerValue) {
		t.Errorf("integer.GetTokenContent not %d, got=%s", integerValue, integer.GetTokenContent())
		return false
	}

	return true
}

//----------------------------------

func TestIdentifierExpressions(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.StatementArray) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.StatementArray))
	}

	stmt, ok := program.StatementArray[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.StatementArray[0] is not ast.ExpressionStatement. got=%T", program.StatementArray[0])
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

func TestIntegerContentExpressions(t *testing.T) {
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
	stmt, ok := program.StatementArray[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.StatementArray[0] is not ast.ExpressionStatement. got=%T", program.StatementArray[0])
	}

	// それの式はIntegerContentだよね？
	integerContent, ok := stmt.Expression.(*ast.Integer)
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		tokentype    token.TokenType
		operator     string
		integerValue int64
	}{
		{"!5", token.BANG, "!", 5},
		{"-15", token.MINUS, "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.StatementArray) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.StatementArray))
		}

		stmt, ok := program.StatementArray[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.StatementArray[0] is not ast.ExpressionStatement. got=%T", program.StatementArray[0])
		}

		if stmt.Token.Type != tt.tokentype {
			t.Fatalf("stmt.Token.Type is not %q. got=%q", tt.tokentype, stmt.Token.Type)
		}

		prefixExp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if prefixExp.Operator != tt.operator {
			t.Fatalf("pex.Operator is not '%s'. got=%s", tt.operator, prefixExp.Operator)
		}
		if !testIntegerContent(t, prefixExp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.StatementArray) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.StatementArray))
		}

		stmt, ok := program.StatementArray[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.StatementArray[0] is not ast.ExpressionStatement. got=%T", program.StatementArray[0])
		}

		expression, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerContent(t, expression.Left, tt.leftValue) {
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not '%s'. got=%s", tt.operator, expression.Operator)
		}

		if !testIntegerContent(t, expression.Right, tt.rightValue) {
			return
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
