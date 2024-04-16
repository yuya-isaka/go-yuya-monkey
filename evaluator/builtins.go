package evaluator

import (
	"fmt"

	"github.com/yuya-isaka/go-yuya-monkey/object"
)

var builtins = map[string]*object.BuiltinObj{
	// 現状ただのラッパー
	"len": &object.BuiltinObj{
		// argsはすでに評価されている（evaluator.go内で）
		// つまりオブジェクトになっている
		Fn: func(args ...object.Object) object.Object {

			if len(args) != 1 {
				return newErrorObj("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			// 1. 文字列
			case *object.StringObj:
				return &object.IntObj{Value: int64(len(arg.Value))}
			// 2. 配列
			case *object.ArrayObj:
				return &object.IntObj{Value: int64(len(arg.Values))}
			default:
				return newErrorObj("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.BuiltinObj{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObj("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newErrorObj("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.ArrayObj)
			if len(array.Values) > 0 {
				return array.Values[0]
			}

			return NULL
		},
	},
	"last": &object.BuiltinObj{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObj("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newErrorObj("argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.ArrayObj)
			length := len(array.Values)
			if length > 0 {
				return array.Values[length-1]
			}

			return NULL
		},
	},
	"rest": &object.BuiltinObj{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObj("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newErrorObj("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.ArrayObj)
			length := len(array.Values)
			if length > 0 {
				newValues := make([]object.Object, length-1)
				copy(newValues, array.Values[1:length])
				return &object.ArrayObj{Values: newValues}
			}

			return NULL
		},
	},
	"push": &object.BuiltinObj{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorObj("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newErrorObj("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			// 1つ目の引数が配列
			array := args[0].(*object.ArrayObj)
			length := len(array.Values)

			newValues := make([]object.Object, length+1)
			copy(newValues, array.Values)
			// 2つ目の引数が追加する値
			newValues[length] = args[1]

			return &object.ArrayObj{Values: newValues}
		},
	},
	"puts": &object.BuiltinObj{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
