package bham

import (
	"fmt"
	"regexp"
	"strings"
	"text/template/parse"
)

var (
	dotVarField = `([\.|\$][^\t^\n^\v^\f^\r^ ]+)+`

	simpleValue    = regexp.MustCompile(`true|false|nil`)
	simpleField    = regexp.MustCompile(fmt.Sprintf(`^%s$`, dotVarField))
	simpleFunction = regexp.MustCompile(fmt.Sprintf(`^([^\.^\t^\n^\v^\f^\r^ ]+)( %s)*$`, dotVarField))
)

func (pt *protoTree) compile() {
	cleanName := pt.name
	i := strings.Index(pt.name, ".bham")
	if i >= 0 {
		cleanName = pt.name[:i] + pt.name[i+5:]
	}

	pt.outputTree = newTree(pt.name, cleanName)

	pt.compileToList(pt.outputTree.Root, pt.nodes)
}

func (pt *protoTree) compileToList(arr *parse.ListNode, nodes []protoNode) {
	for _, node := range nodes {
		switch node.identifier {
		case identRaw:
			arr.Nodes = append(arr.Nodes, newTextNode(node.content))
		case identFilter:
			if node.needsRuntimeData() {
			} else {
				content := node.filter.Open + node.filter.Handler(node.content) + node.filter.Close
				arr.Nodes = append(arr.Nodes, newTextNode(content))
			}
		case identExecutable:
			switch {
			case simpleValue.MatchString(node.content):
				arr.Nodes = append(arr.Nodes, newValueNode(node.content))
			case simpleField.MatchString(node.content):
				arr.Nodes = append(arr.Nodes, newFieldNode(node.content))
			case simpleFunction.MatchString(node.content):
				arr.Nodes = append(arr.Nodes, newFunctionNode(node.content))
			default:
				node, err := processCode(node.content)
				if err == nil {
					arr.Nodes = append(arr.Nodes, node)
				} else {
					pt.err = err
					return
				}
			}
		case identTag:
			nodes, err := newStandaloneTag(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, nodes...)
			} else {
				pt.err = err
			}
		case identTagOpen:
			td, _, err := parseTag(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, newTextNode(td.Opening()))
			} else {
				pt.err = err
			}
		case identTagClose:
			td, _, err := parseTag(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, newTextNode(td.Close()))
			} else {
				pt.err = err
			}
		case identText:
			arr.Nodes = append(arr.Nodes, newMaybeTextNode(node.content)...)
		case identIf:
			var err error
			branching := &parse.PipeNode{
				NodeType: parse.NodePipe,
				Cmds: []*parse.CommandNode{
					&parse.CommandNode{
						NodeType: parse.NodeCommand,
						Args:     []parse.Node{},
					},
				},
			}
			switch {
			case simpleValue.MatchString(node.content):
				branching.Cmds[0].Args = append(
					branching.Cmds[0].Args,
					newValueNode(node.content),
				)
			case simpleField.MatchString(node.content):
				branching.Cmds[0].Args = append(
					branching.Cmds[0].Args,
					newFieldNode(node.content),
				)
			case simpleFunction.MatchString(node.content):
				branching = newFieldNode(node.content).(*parse.ActionNode).Pipe
			default:
				node, e := processCode(node.content)
				err = e
				if err == nil {
					branching = node.Pipe
				} else {
					pt.err = err
					return
				}
			}

			if err == nil {
				in := &parse.IfNode{
					parse.BranchNode{
						NodeType: parse.NodeIf,
						Pipe:     branching,
						List: &parse.ListNode{
							NodeType: parse.NodeList,
						},
						ElseList: &parse.ListNode{
							NodeType: parse.NodeList,
						},
					},
				}
				if len(node.list) > 0 {
					pt.compileToList(in.List, node.list)
				}
				if len(node.elseList) > 0 {
					pt.compileToList(in.ElseList, node.elseList)
				}
				arr.Nodes = append(arr.Nodes, in)
			} else {
				pt.err = err
			}
		default:
			fmt.Println(node.identifier)
			fmt.Println(node.content)
		}
	}
}
