package evaluator

import (
	"github.com/yuya-isaka/go-yuya-monkey/ast"
	"github.com/yuya-isaka/go-yuya-monkey/object"
)

func quote(node ast.Node) object.Object {
	return &object.QuoteObj{Node: node}
}
