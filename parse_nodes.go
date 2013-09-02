package bham

import (
	"strings"
	"text/template/parse"
)

func newTree(source, name string) *parse.Tree {
	return &parse.Tree{
		Name:      name,
		ParseName: source,
		Root: &parse.ListNode{
			NodeType: parse.NodeList,
		},
	}
}
func newTextNode(text string) parse.Node {
	return &parse.TextNode{
		NodeType: parse.NodeText,
		Text:     []byte(text),
	}
}

func newFieldNode(field string) parse.Node {
	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						&parse.FieldNode{
							NodeType: parse.NodeField,
							Ident:    strings.Split(field, "."),
						},
					},
				},
			},
		},
	}
}
