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
		{
			`{"name": "Monkey"}[fn(x) {x}];`,
			"unusable as hash key: FUNCTION",
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
	str, ok := obj.(*object.StringObj)
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
	str, ok := obj.(*object.StringObj)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", obj, obj)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INT"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, nil},
		{`first([1, 2, 3])`, 1},
		{`first(1)`, "argument to `first` must be ARRAY, got INT"},
		{`last([1, 2, 3])`, 3},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "argument to `push` must be ARRAY, got INT"},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)

		switch expect := tt.expect.(type) {
		case int:
			testIntObj(t, obj, int64(expect))
		case nil:
			testNullObj(t, obj)
		case string:
			errObj, ok := obj.(*object.ErrorObj)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
				// 続けないと終わっちゃうので
				continue
			}

			if errObj.Value != expect {
				t.Errorf("wrong error message. expected=%q, got=%q", expect, errObj.Value)
			}
		case []int:
			array, ok := obj.(*object.ArrayObj)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", obj, obj)
				continue
			}

			if len(array.Values) != len(expect) {
				t.Errorf("wrong num of values. want=%d, got=%d", len(expect), len(array.Values))
				continue
			}

			// 一旦今回のテストはint64に絞っている
			// 本来はどの型でも配列は持てる（テストしていないだけ）
			for i, expectValue := range expect {
				testIntObj(t, array.Values[i], int64(expectValue))
			}
		}
	}
}

func TestArray(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	obj := testEval(input)
	result, ok := obj.(*object.ArrayObj)
	if !ok {
		t.Fatalf("obj is not ArrayObj. got=%T(%+v)", obj, obj)
	}

	if len(result.Values) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Values))
	}

	testIntObj(t, result.Values[0], 1)
	testIntObj(t, result.Values[1], 4)
	testIntObj(t, result.Values[2], 6)
}

func TestArrayIndex(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{
			"[1,2,3][0]",
			1,
		},
		{
			"[1,2,3][1]",
			2,
		},
		{
			"[1,2,3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1,2,3][1+1];",
			3,
		},
		{
			"let myArray = [1,2,3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1,2,3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1,2,3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1,2,3][3]",
			nil,
		},
		{
			"[1,2,3][-1]",
			nil,
		},
		{
			"let three = fn(x) { return [x,x+1,x+2]; }; three(2)[1]",
			3,
		},
		{
			"let three = fn(x) { return [x,x+1,x+2]; }; three(10)[0]",
			10,
		},
		{
			"let three = fn(x) { return [x,x+1,x+2]; }; three(20)[2]",
			22,
		},
		{
			"-[2,3,4][2]",
			-4,
		},
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

func TestHash(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6/ 2,
		4: 4,
		true: 5,
		false: 6
	}`

	obj := testEval(input)
	result, ok := obj.(*object.HashObj)
	if !ok {
		t.Fatalf("Eval didn't return HashObj. got=%T (%+v)", obj, obj)
	}

	expect := map[object.HashKey]int64{
		(&object.StringObj{Value: "one"}).HashKey():   1,
		(&object.StringObj{Value: "two"}).HashKey():   2,
		(&object.StringObj{Value: "three"}).HashKey(): 3,
		(&object.IntObj{Value: 4}).HashKey():          4,
		TRUE.HashKey():                                5,
		FALSE.HashKey():                               6,
	}

	if len(result.Pairs) != len(expect) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectKey, expectValue := range expect {
		pair, ok := result.Pairs[expectKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntObj(t, pair.Value, expectValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntObj(t, evaluated, int64(integer))
		} else {
			testNullObj(t, evaluated)
		}
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
