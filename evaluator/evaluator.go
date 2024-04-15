package evaluator

// 評価とは、Nodeインタフェースを満たした構造体をObjectインタフェースを満たした構造体にすること

// このフェーズで、初めて、プログラミング言語に意味があらわれる

import (
	"fmt"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/object"
)

// 勿体無いから事前に生成して、参照する
var (
	NULL  = &object.NullObj{}
	TRUE  = &object.BoolObj{Value: true}
	FALSE = &object.BoolObj{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {

	// [ノード → オブジェクト]しまくる
	//		- 文はそのままEval呼び出し
	// 		- 式はそのままEval呼び出し
	// 		- 他のやつは変換
	//    - リターンだったら即返却
	switch node := node.(type) {
	case *ast.ProgramNode:
		var obj object.Object

		for _, statement := range node.Statements {
			obj = Eval(statement, env)

			// 中身を取り出すには型アサーション必要
			switch obj := obj.(type) {
			// リターンなら中身取り出し
			case *object.ReturnObj:
				return obj.Value
			// エラーならそのまま
			case *object.ErrorObj:
				return obj
			}
		}

		// 最後に評価した結果を返す
		// Returnあったらそれを事前に返している
		return obj

	case *ast.BlockNode:
		var obj object.Object

		for _, statement := range node.Statements {
			obj = Eval(statement, env)

			if obj != nil {
				vt := obj.Type()
				if vt == object.RETURN || vt == object.ERROR {
					// そのまま上に上げる
					// ProgramNodeのところで取り出すため (ProgramNodeのところで終われない)
					return obj
				}
			}
		}

		// 最後に評価した結果を返す
		// Returnあったらそれを事前に返している
		return obj

	case *ast.LetNode:
		obj := Eval(node.Value, env)
		if isErrorObj(obj) {
			return obj
		}

		// 環境に登録
		env.Set(node.Name.Value, obj)

	case *ast.ReturnNode:
		obj := Eval(node.Value, env)
		if isErrorObj(obj) {
			return obj
		}
		return &object.ReturnObj{Value: obj}

	case *ast.EsNode:
		return Eval(node.Value, env)

	case *ast.IdentNode:
		if obj, ok := env.Get(node.Value); ok {
			return obj
		}
		// 環境から見つけられなかったらフォールバックして組込み関数調べる
		if builtin, ok := builtins[node.Value]; ok {
			return builtin
		}

		return newErrorObj("identifier not found: " + node.Value)

	case *ast.IntNode:
		return &object.IntObj{Value: node.Value}

	case *ast.BoolNode:
		return changeBoolObj(node.Value)

	case *ast.StringNode:
		return &object.StringObj{Value: node.Value}

	case *ast.PrefixNode:
		right := Eval(node.Right, env)
		if isErrorObj(right) {
			return right
		}
		return evalPrefix(node.Operator, right)

	case *ast.InfixNode:
		left := Eval(node.Left, env)
		if isErrorObj(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isErrorObj(right) {
			return right
		}

		switch {

		// 1.
		// 先にInt型をやることで、「先に==や!=で変換して比較される」ことを防いでいる
		case left.Type() == object.INT && right.Type() == object.INT:

			// 値をオブジェクトからアンラップ
			leftVal := left.(*object.IntObj).Value
			rightVal := right.(*object.IntObj).Value

			switch node.Operator {
			case "+":
				return &object.IntObj{Value: leftVal + rightVal}
			case "-":
				return &object.IntObj{Value: leftVal - rightVal}
			case "*":
				return &object.IntObj{Value: leftVal * rightVal}
			case "/":
				return &object.IntObj{Value: leftVal / rightVal}
			case "<":
				return changeBoolObj(leftVal < rightVal)
			case ">":
				return changeBoolObj(leftVal > rightVal)
			case "==":
				return changeBoolObj(leftVal == rightVal)
			case "!=":
				return changeBoolObj(leftVal != rightVal)
			default:
				return newErrorObj("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
			}

		// 2.
		case left.Type() == object.STRING && right.Type() == object.STRING:
			if node.Operator != "+" && node.Operator != "==" && node.Operator != "!=" {
				return newErrorObj("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
			}

			leftVal := left.(*object.StringObj).Value
			rightVal := right.(*object.StringObj).Value
			switch node.Operator {
			case "+":
				return &object.StringObj{Value: leftVal + rightVal}
			case "==":
				return changeBoolObj(leftVal == rightVal)
			case "!=":
				return changeBoolObj(leftVal != rightVal)
			}

		// 3.
		// オブジェクトを指し示すのにポインタ（参照）のみを使っていて、ポインタを比較すればいい
		// 		ポインタ（配置されているメモリアドレス）を比較している
		//  	オブジェクトは、整数かTRUEかFALSEかNULLだけ。整数は先に計算して、残りは参照だけ
		//		整数や他のオブジェクトはポインタの比較を単純にするわけにはいかない（毎回新しく生成しているから）
		case node.Operator == "==":
			return changeBoolObj(left == right)
		case node.Operator == "!=":
			return changeBoolObj(left != right)

		// 4.
		case left.Type() != right.Type():
			return newErrorObj("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())

		default:
			return newErrorObj("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
		}

	case *ast.IfNode:
		condition := Eval(node.Condition, env)
		if isErrorObj(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(node.Consequence, env)
		} else if node.Alternative != nil {
			return Eval(node.Alternative, env)
		} else {
			return NULL
		}

	case *ast.FunctionNode:
		// けっこうそのままいれる
		return &object.FunctionObj{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallNode:
		function := Eval(node.Function, env)
		if isErrorObj(function) {
			return function
		}

		// これの方が効率がいい
		args := make([]object.Object, len(node.Arguments))
		// 引数を左から右に評価
		for i, e := range node.Arguments {
			obj := Eval(e, env)
			if isErrorObj(obj) {
				return obj
			}
			args[i] = obj
		}

		switch fn := function.(type) {
		case *object.FunctionObj:
			// パラメータを『拡張した環境』に束縛
			extendedEnv := object.NewEnclosedEnvironment(fn.Env)
			for paramIdx, param := range fn.Parameters {
				// パラメータの変数 ← 評価結果
				env.Set(param.Value, args[paramIdx])
			}

			// ボディと『拡張した環境』で評価
			result := Eval(fn.Body, extendedEnv)

			// もし、結果がReturnオブジェクトだったらそのまま返却
			// その関数からのリターンだから、これはBlockの時みたいに上に上げなくていい
			// むしろこのif文がないと、そのままReturnが浮上して処理が止まってしまう
			if returnValue, ok := result.(*object.ReturnObj); ok {
				return returnValue.Value
			}

			return result

		// 絶対にReturnを返さないので、アンラップする必要がない
		case *object.BuiltinObj:
			// 引数の評価結果をそのまま渡す
			// builtins.go内でよしなに処理
			return fn.Fn(args...)

		default:
			return newErrorObj("not a function: %s", function.Type()) // 存在していないfn.Type()しててランタイムエラーになっていた
		}

	case *ast.ArrayNode:
		values := make([]object.Object, len(node.Values))
		for i, v := range node.Values {
			obj := Eval(v, env)
			if isErrorObj(obj) {
				return obj
			}
			values[i] = obj
		}

		return &object.ArrayObj{Values: values}

	case *ast.IndexNode:
		left := Eval(node.Left, env)
		if isErrorObj(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isErrorObj(index) {
			return index
		}

		switch {
		case left.Type() == object.ARRAY && index.Type() == object.INT:
			array := left.(*object.ArrayObj)
			idx := index.(*object.IntObj).Value
			max := int64(len(array.Values) - 1)

			if idx < 0 || max < idx {
				return NULL
			}

			return array.Values[idx]

		default:
			return newErrorObj("index operator not supported: %s", left.Type())
		}

	}

	return nil
}

func changeBoolObj(value bool) object.Object {
	if value {
		return TRUE
	} else {
		return FALSE
	}
}

// nullでもなく、falseでもないやつ
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// ...--で受け取り
//
//		--...で渡す
//	  ーーなんだよね〜....みたいなイメージかな（違う）
func newErrorObj(format string, a ...interface{}) *object.ErrorObj {
	return &object.ErrorObj{Value: fmt.Sprintf(format, a...)}
}

func isErrorObj(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {

	// 否定
	case "!":
		switch right {
		case TRUE:
			return FALSE
		case FALSE:
			return TRUE
		case NULL:
			return TRUE
		default:
			return FALSE
		}

	// マイナス
	case "-":
		if right.Type() != object.INT {
			return newErrorObj("unknown operator: -%s", right.Type())
		}

		value := right.(*object.IntObj).Value
		return &object.IntObj{Value: -value}

	default:
		return newErrorObj("unknown operator: %s%s", operator, right.Type())
	}

}
