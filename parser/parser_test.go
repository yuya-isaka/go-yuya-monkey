package parser

import (
	"fmt"
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/lexer"
)

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

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectIdent string
		expectValue interface{}
	}{
		// {"let x =    5      a  ;", "x", 5}, // エラーになるよ
		{"let x =    5        ;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// 1. 長さチェック
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		// 2. 文自体チェック
		stmt := program.Statements[0]
		if !testLetStatementIs(t, stmt, tt.expectIdent) {
			return
		}

		// 3. 中身の式チェック
		expression := stmt.(*ast.LetNode).Value
		if !testContentExpression(t, expression, tt.expectValue) {
			return
		}
	}

}

//----------------------------------

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnNode)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ReturnStatement. got=%T", program.Statements[0])
		}

		if returnStmt.Token.Name != "return" {
			t.Fatalf("returnStmt.GetTokenContent not 'return', got %q", returnStmt.Token.Name)
		}
		if testContentExpression(t, returnStmt.Value, tt.expectValue) {
			return
		}
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Value.(*ast.IdentNode)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Value)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.Token.Name != "foobar" {
		t.Errorf("ident.GetTokenContent not %s. got=%s", "foobar", ident.Token.Name)
	}
}

func TestIntegerContentExpressions(t *testing.T) {
	input := "42;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// 正しい数の文が生成されたか
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	// 最初の文がExpressionStatementだよね？
	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// それの式はIntegerContentだよね？
	integerContent, ok := stmt.Value.(*ast.IntNode)
	if !ok {
		t.Fatalf("exp not *ast.IntegerContent. got=%T", stmt.Value)
	}
	// それの中身は42だよね？
	if integerContent.Value != 42 {
		t.Errorf("integerContent.IntegerValue not %d. got=%d", 42, integerContent.Value)
	}
	// トークンとして取得したら"42"だよね？
	if integerContent.Token.Name != "42" {
		t.Errorf("integerContent.GetTokenContent not %s. got=%s", "42", integerContent.Token.Name)
	}

	// 上記のアサーションを設ける
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.EsNode)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := stmt.Value.(*ast.PrefixNode)
		if !ok {
			t.Fatalf("PrefixExpressionへの変換失敗！ got=%T", stmt.Value)
		}
		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not %s. got=%s", tt.operator, expression.Operator)
		}
		if !testContentExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.EsNode)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Value, tt.leftValue, tt.operator, tt.rightValue) {
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
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
		{
			"add(a, b)[3]",
			"(add(a, b)[3])",
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

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("長さ違うよ got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.EsNode)
		if !ok {
			t.Fatalf("入ってるのが式文じゃない got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Value.(*ast.BoolNode)
		if !ok {
			t.Fatalf("式がブーリアンじゃないな got=%T", stmt.Value)
		}

		if boolean.Value != tt.expectBoolean {
			t.Errorf("思ったやつじゃない%tが欲しいが got=%t", tt.expectBoolean, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := stmt.Value.(*ast.IfNode)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Value)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Value, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative.Statements was not nil. got=%+v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := stmt.Value.(*ast.IfNode)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Value)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Value, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d\n", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Value, "y") {
		return
	}
}

func TestFunctionParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Value.(*ast.FunctionNode)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Funciton. got=%T", stmt.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function parameterList wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testContentExpression(t, function.Parameters[0], "x")
	testContentExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.StatementArray has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("function.Body.StatmentList[0] is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Value, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input           string
		expectParamList []string
	}{
		{input: "fn() {};", expectParamList: []string{}},
		{input: "fn(x) {};", expectParamList: []string{"x"}},
		{input: "fn(x, y, z) {};", expectParamList: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.EsNode)
		function := stmt.Value.(*ast.FunctionNode)

		if len(function.Parameters) != len(tt.expectParamList) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectParamList), len(function.Parameters))
		}

		for i, ident := range tt.expectParamList {
			testContentExpression(t, function.Parameters[i], ident)
		}
	}

}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statement. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := stmt.Value.(*ast.CallNode)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Value)
	}

	if !testIdentifier(t, expression.Function, "add") {
		return
	}

	if len(expression.Arguments) != 3 {
		t.Fatalf("wrong length of argumetlist. got=%d", len(expression.Arguments))
	}

	testContentExpression(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func TestCallExpressinParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectIdent   string
		expectArgList []string
	}{
		{
			input:         "add();",
			expectIdent:   "add",
			expectArgList: []string{},
		},
		{
			input:         "add(1)",
			expectIdent:   "add",
			expectArgList: []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectIdent:   "add",
			expectArgList: []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.EsNode)
		expression := stmt.Value.(*ast.CallNode)

		if !testIdentifier(t, expression.Function, tt.expectIdent) {
			return
		}

		if len(expression.Arguments) != len(tt.expectArgList) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d\n", len(tt.expectArgList), len(expression.Arguments))
		}

		for i, arg := range tt.expectArgList {
			if expression.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i, arg, expression.Arguments[i].String())
			}
		}
	}
}

func TestString(t *testing.T) {
	input := `"hello world";`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.EsNode)
	node, ok := stmt.Value.(*ast.StringNode)
	if !ok {
		t.Fatalf("exp not *ast.StringNode. got=%T", stmt.Value)
	}

	if node.Value != "hello world" {
		t.Errorf("node.Value not %q. got=%q", "hello world", node.Value)
	}
}

func TestParsingArray(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("exp not ast.EsNode. got=%T", program.Statements[0])
	}
	array, ok := stmt.Value.(*ast.ArrayNode)
	if !ok {
		t.Fatalf("exp not ast.ArrayNode. got=%T", stmt.Value)
	}

	if len(array.Values) != 3 {
		t.Fatalf("len(array.Values) not 3. got=%d", len(array.Values))
	}

	testIntegerContent(t, array.Values[0], 1)
	testInfixExpression(t, array.Values[1], 2, "*", 2)
	testInfixExpression(t, array.Values[2], 3, "+", 3)
}

func TestParsingIndex(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.EsNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.EsNode. got=%T", program.Statements[0])
	}
	index, ok := stmt.Value.(*ast.IndexNode)
	if !ok {
		t.Fatalf("stmt.Value is not ast.IndexNode. got=%T", stmt.Value)
	}

	if !testIdentifier(t, index.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, index.Index, 1, "+", 1) {
		return
	}
}

// -------------------------------- ヘルパー関数

// テスト、文、期待値
func testLetStatementIs(t *testing.T, statement ast.Statement, expectName string) bool {
	letStmt, ok := statement.(*ast.LetNode)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStmt.Name.Value != expectName {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", expectName, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.Token.Name != expectName {
		t.Errorf("letStmt.Name.TokenContent() not '%s'. got=%s", expectName, letStmt.Name.Token.Name)
		return false
	}

	return true
}

// 整数
// テスト、式、期待値
func testIntegerContent(t *testing.T, expression ast.Expression, expectValue int64) bool {
	// 1. 型変換
	integer, ok := expression.(*ast.IntNode)
	if !ok {
		t.Errorf("expression not *ast.Integer. got=%T", expression)
		return false
	}

	// 2. 期待した値を持っているか
	if integer.Value != expectValue {
		t.Errorf("integer.IntegerValue not %d, got=%d", expectValue, integer.Value)
		return false
	}

	// 3. トークン確認
	// 文字列にしてから比較しようね〜
	if integer.Token.Name != fmt.Sprintf("%d", expectValue) {
		t.Errorf("integer.GetTokenContent not %d, got=%s", expectValue, integer.Token.Name)
		return false
	}

	return true
}

// 識別子
// テスト、式、期待値
func testIdentifier(t *testing.T, expression ast.Expression, expectValue string) bool {

	// 1. 型変換
	identifier, ok := expression.(*ast.IdentNode)
	if !ok {
		t.Errorf("expression not *ast.Identifier. got=%T", expression)
		return false
	}

	// 2. 期待した値を持っているか
	if identifier.Value != expectValue {
		t.Errorf("ident.IdentValue not %s. got=%s", expectValue, identifier.Value)
		return false
	}

	// 3. トークン確認
	if identifier.Token.Name != expectValue {
		t.Errorf("ident.GetTokenContent not %s. got=%s", expectValue, identifier.Token.Name)
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
	case bool:
		return testBooleanContent(t, expression, v)
	}

	t.Errorf("type of exp not handled. got=%T", expression)
	return false
}

// テスト、式、左辺、オペレータ、右辺
func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {

	// 1. 型変換
	infixExpression, ok := expression.(*ast.InfixNode)
	if !ok {
		t.Errorf("expression is not ast.InfixExpression. got=%T(%s)", expression, expression)
		return false
	}

	// 左辺チェック
	if !testContentExpression(t, infixExpression.Left, left) {
		return false
	}

	// 文字列チェック
	if infixExpression.Operator != operator {
		t.Errorf("expression.Operator is not '%s'. got=%q", operator, infixExpression.Operator)
		return false
	}

	if !testContentExpression(t, infixExpression.Right, right) {
		return false
	}

	return true
}

func testBooleanContent(t *testing.T, expression ast.Expression, expectValue bool) bool {
	boolean, ok := expression.(*ast.BoolNode)
	if !ok {
		t.Errorf("expression not *ast.Boolean. got=%T", expression)
		return false
	}

	if boolean.Value != expectValue {
		t.Errorf("boolean.BoolValue not %t. got=%t", expectValue, boolean.Value)
		return false
	}

	if boolean.Token.Name != fmt.Sprintf("%t", expectValue) {
		t.Errorf("boolean.GetTokenContent not %t. got=%s", expectValue, boolean.Token.Name)
		return false
	}

	return true
}
