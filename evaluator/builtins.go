package evaluator

import "github.com/yuya-isaka/go-yuya-monkey/object"

var builtins = map[string]*object.BuiltinObj{
	// 現状ただのラッパー
	"len": &object.BuiltinObj{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObj("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.StringObj:
				return &object.IntObj{Value: int64(len(arg.Value))}
			default:
				return newErrorObj("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
