package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntNode{Value: 1} }
	two := func() Expression { return &IntNode{Value: 2} }

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntNode)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input  Node
		expect Node
	}{
		{
			one(),
			two(),
		},
		{
			&ProgramNode{
				Statements: []Statement{
					&EsNode{Value: one()},
				},
			},
			&ProgramNode{
				Statements: []Statement{
					&EsNode{Value: two()},
				},
			},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expect)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expect)
		}
	}
}
