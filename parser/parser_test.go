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
		// {
		// 	"a + b + c",
		// 	"((a + b) + c)",
		// },
		// {
		// 	"a + b - c",
		// 	"((a + b) - c)",
		// },
		// {
		// 	"a * b * c",
		// 	"((a * b) * c)",
		// },
		// {
		// 	"a * b / c",
		// 	"((a * b) / c)",
		// },
		// {
		// 	"a + b / c",
		// 	"(a + (b / c))",
		// },
		// {
		// 	"a + b * c + d / e - f",
		// 	"(((a + (b * c)) + (d / e)) - f)",
		// },
		// {
		// 	"3 + 4; -5 * 5",
		// 	"(3 + 4)((-5) * 5)",
		// },
		// {
		// 	"5 > 4 == 3 < 4",
		// 	"((5 > 4) == (3 < 4))",
		// },
		// {
		// 	"5 < 4 != 3 > 4",
		// 	"((5 < 4) != (3 > 4))",
		// },
		// {
		// 	"3 + 4 * 5 == 3 * 1 + 4 * 5",
		// 	"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		// },
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

// -------------------------------- ヘルパー関数

// テスト、パーサー
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	// ok
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

// テスト、文、期待値
func testLetStatementIs(t *testing.T, statement ast.Statement, expectName string) bool {
	if statement.GetTokenContent() != "let" {
		t.Errorf("statement.TokenContent not 'let'. got=%q", statement.GetTokenContent())
		return false
	}

	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStmt.IdentName.IdentValue != expectName {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", expectName, letStmt.IdentName.IdentValue)
		return false
	}

	if letStmt.IdentName.GetTokenContent() != expectName {
		t.Errorf("letStmt.Name.TokenContent() not '%s'. got=%s", expectName, letStmt.IdentName.GetTokenContent())
		return false
	}

	return true
}

// 整数
// テスト、式、期待値
func testIntegerContent(t *testing.T, expression ast.Expression, expectValue int64) bool {
	integer, ok := expression.(*ast.Integer)
	if !ok {
		t.Errorf("expression not *ast.Integer. got=%T", expression)
		return false
	}

	if integer.IntegerValue != expectValue {
		t.Errorf("integer.IntegerValue not %d, got=%d", expectValue, integer.IntegerValue)
		return false
	}

	// 文字列にしてから比較しようね〜
	if integer.GetTokenContent() != fmt.Sprintf("%d", expectValue) {
		t.Errorf("integer.GetTokenContent not %d, got=%s", expectValue, integer.GetTokenContent())
		return false
	}

	return true
}

// 識別子
// テスト、式、期待値
func testIdentifier(t *testing.T, expression ast.Expression, expectValue string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression not *ast.Identifier. got=%T", expression)
		return false
	}

	if identifier.IdentValue != expectValue {
		t.Errorf("ident.IdentValue not %s. got=%s", expectValue, identifier.IdentValue)
		return false
	}

	if identifier.GetTokenContent() != expectValue {
		t.Errorf("ident.GetTokenContent not %s. got=%s", expectValue, identifier.GetTokenContent())
		return false
	}

	return true
}

func testContentExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerContent(t, expression, int64(v))
	case int64:
		return testIntegerContent(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	}

	t.Errorf("type of exp not handled. got=%T", expression)
	return false
}

// テスト、式、左辺、オペレータ、右辺
func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.InfixExpression. got=%T(%s)", expression, expression)
		return false
	}

	if !testContentExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("expression.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testContentExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
