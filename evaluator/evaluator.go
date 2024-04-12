package evaluator

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
	switch node := node.(type) {
	case *ast.ProgramNode:
		var obj object.Object

		for _, statement := range node.Statements {
			obj = Eval(statement, env)

			// 中身を取り出すには型アサーション必要
			switch obj := obj.(type) {
			case *object.ReturnObj:
				return obj.Value
			case *object.ErrorObj:
				return obj
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

		env.Set(node.Name.Value, obj)

	case *ast.IdentNode:
		obj, ok := env.Get(node.Value)
		if !ok {
			return newErrorObj("identifier not found: " + node.Value)
		}

		return obj

	case *ast.EsNode:
		return Eval(node.Value, env)

	case *ast.IntNode:
		return &object.IntObj{Value: node.Value}

	case *ast.BoolNode:
		return changeBoolObj(node.Value)

	case *ast.StringNode:
		return &object.String{Value: node.Value}

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

		// オブジェクトを指し示すのにポインタ（参照）のみを使っていて、ポインタを比較すればいい
		// 		ポインタ（配置されているメモリアドレス）を比較している
		//  	オブジェクトは、整数かTRUEかFALSEかNULLだけ。整数は先に計算して、残りは参照だけ
		//		整数や他のオブジェクトはポインタの比較を単純にするわけにはいかない（毎回新しく生成しているから）
		case node.Operator == "==":
			return changeBoolObj(left == right)
		case node.Operator == "!=":
			return changeBoolObj(left != right)

		case left.Type() != right.Type():
			return newErrorObj("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())

		default:
			return newErrorObj("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
		}

	case *ast.BlockNode:
		var obj object.Object

		for _, statement := range node.Statements {
			obj = Eval(statement, env)

			if obj != nil {
				vt := obj.Type()
				if vt == object.RETURN || vt == object.ERROR {
					// そのまま上に上げる
					return obj
				}
			}
		}

		return obj

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

	case *ast.ReturnNode:
		obj := Eval(node.Value, env)
		if isErrorObj(obj) {
			return obj
		}
		return &object.ReturnObj{Value: obj}

	case *ast.FunctionNode:
		return &object.FunctionObj{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallNode:
		function := Eval(node.Function, env)
		if isErrorObj(function) {
			return function
		}

		var args []object.Object
		// 引数を左から右に評価
		for _, e := range node.Arguments {
			obj := Eval(e, env)
			if isErrorObj(obj) {
				return obj
			}
			args = append(args, obj)
		}

		fn, ok := function.(*object.FunctionObj)
		if !ok {
			return newErrorObj("not a function: %s", fn.Type())
		}

		// パラメータを『拡張した環境』に束縛
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		for paramIdx, param := range fn.Parameters {
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
