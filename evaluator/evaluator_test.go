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
		obj := testEval(tt.input)
		testIntObj(t, obj, tt.expect)
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
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBoolObj(t, obj, tt.expect)
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
		obj := testEval(tt.input)
		testBoolObj(t, obj, tt.expect)
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
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntObj(t, obj, tt.expect)
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

// --------------------------------

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	return Eval(program)
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
