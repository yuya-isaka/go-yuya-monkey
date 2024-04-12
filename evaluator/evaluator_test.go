package evaluator

import (
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/lexer"
	"github.com/yuya-isaka/go-yuya-monkey/object"
	"github.com/yuya-isaka/go-yuya-monkey/parser"
)

func TestEvalInt(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"5", 5},
		{"10", 10},
		{"10;", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		testIntObj(t, testEval(tt.input), tt.expect)
	}
}

func TestEvalBool(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"foobar" == "foobar"`, true},
		{`"foobr" == "foobar"`, false},
		{`"foobr" != "foobar"`, true},
		{`"foobar" != "foobar"`, false},
	}

	for _, tt := range tests {
		testBoolObj(t, testEval(tt.input), tt.expect)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		testBoolObj(t, testEval(tt.input), tt.expect)
	}
}

func TestIfElse(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		integer, ok := tt.expect.(int)
		if ok {
			testIntObj(t, obj, int64(integer))
		} else {
			testNullObj(t, obj)
		}
	}
}

func TestReturn(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`, 10,
		},
		{
			`
			let f = fn(x) {
				return x;
				x + 10;
			};
			f(10);
			`, 10,
		},
		{
			`
			let f = fn(x) {
				let result = x + 10;
				return result;
				return 10;
			};
			f(10);
			`, 20,
		},
		{
			`
			let f = fn(x) {
				let result = x + 10;
				return result;
				return 10;
			};
			f(10);
			return 30;
			`, 30,
		},
	}

	for _, tt := range tests {
		testIntObj(t, testEval(tt.input), tt.expect)
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			"5 + true;",
			"type mismatch: INT + BOOL",
		},
		{
			"5 + true; 5;",
			"type mismatch: INT + BOOL",
		},
		{
			"-true",
			"unknown operator: -BOOL",
		},
		{
			"true + false;",
			"unknown operator: BOOL + BOOL",
		},
		{
			"5; true + false; 10;",
			"unknown operator: BOOL + BOOL",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOL + BOOL",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOL + BOOL",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)

		errObj, ok := obj.(*object.ErrorObj)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", obj, obj)
			continue
		}

		if errObj.Value != tt.expect {
			t.Errorf("wrong error message. expect=%q, got=%q", tt.expect, errObj.Value)
		}
	}
}

func TestLet(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntObj(t, testEval(tt.input), tt.expect)
	}
}

func TestFunction(t *testing.T) {
	input := "fn(x) { x + 2 };"
	obj := testEval(input)

	fn, ok := obj.(*object.FunctionObj)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", obj, obj)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectBody := "(x + 2)"
	if fn.Body.String() != expectBody {
		t.Fatalf("body is not %q. got=%q", expectBody, fn.Body.String())
	}
}

func TestCall(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let ident = fn(x) { x; }; ident(5);", 5},
		{"let ident = fn(x) { return x; }; ident(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntObj(t, testEval(tt.input), tt.expect)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntObj(t, testEval(input), 4)
}

func TestString(t *testing.T) {
	input := `"Hello World!"`

	obj := testEval(input)
	str, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", obj, obj)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	obj := testEval(input)
	str, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", obj, obj)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

// --------------------------------

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

// --------------------------------

func testIntObj(t *testing.T, obj object.Object, expect int64) bool {
	result, ok := obj.(*object.IntObj)
	if !ok {
		t.Errorf("object is not IntObj. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expect {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expect)
		return false
	}

	return true
}

func testBoolObj(t *testing.T, obj object.Object, expect bool) bool {
	result, ok := obj.(*object.BoolObj)
	if !ok {
		t.Errorf("object is not BoolObj. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expect {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expect)
		return false
	}

	return true
}

func testNullObj(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
