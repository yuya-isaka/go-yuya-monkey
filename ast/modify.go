package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {

	switch node := node.(type) {

	case *ProgramNode:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *EsNode:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	}

	return modifier(node)
}
