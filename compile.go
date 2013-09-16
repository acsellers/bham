package bham

import (
	"fmt"
	"regexp"
	"strings"
	"text/template/parse"
)

var (
	dotVarField = `([\.|\$][^\t^\n^\v^\f^\r^ ]+)+`

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
		default:
			fmt.Println(node.identifier)
			fmt.Println(node.content)
		}
	}
}
