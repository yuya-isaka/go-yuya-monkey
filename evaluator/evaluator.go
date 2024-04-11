package evaluator

import (
	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/object"
)

// 勿体無いから事前に生成して、参照する
var (
	NULL  = &object.NullObj{}
	TRUE  = &object.BoolObj{Value: true}
	FALSE = &object.BoolObj{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.ProgramNode:
		return evalStatements(node.Statements)

	case *ast.EsNode:
		return Eval(node.Value)

	case *ast.IntNode:
		return &object.IntObj{Value: node.Value}

	case *ast.BoolNode:
		return changeBoolObj(node.Value)

	case *ast.PrefixNode:
		right := Eval(node.Right)
		return evalPrefix(node.Operator, right)

	case *ast.InfixNode:
		left := Eval(node.Left)
		right := Eval(node.Right)

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
				return NULL
			}

		// オブジェクトを指し示すのにポインタ（参照）のみを使っていて、ポインタを比較すればいい
		// 		ポインタ（配置されているメモリアドレス）を比較している
		//  	オブジェクトは、整数かTRUEかFALSEかNULLだけ。整数は先に計算して、残りは参照だけ
		//		整数や他のオブジェクトはポインタの比較を単純にするわけにはいかない（毎回新しく生成しているから）
		case node.Operator == "==":
			return changeBoolObj(left == right)
		case node.Operator == "!=":
			return changeBoolObj(left != right)

		default:
			return NULL
		}

	case *ast.BlockNode:
		return evalStatements(node.Statements)

	case *ast.IfNode:
		condition := Eval(node.Condition)

		if isTruthy(condition) {
			return Eval(node.Consequence)
		} else if node.Alternative != nil {
			return Eval(node.Alternative)
		} else {
			return NULL
		}
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	// 最後に評価した結果を返す
	return result
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
			return NULL
		}

		value := right.(*object.IntObj).Value
		return &object.IntObj{Value: -value}

	default:
		return NULL
	}

}
