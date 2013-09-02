package bham

import (
	"strings"
	"text/template/parse"
)

func (pt *protoTree) compile() {
	cleanName := pt.name
	i := strings.Index(pt.name, ".bham")
	if i >= 0 {
		cleanName = pt.name[:i] + pt.name[i+5:]
	}

	pt.outputTree = newTree(pt.name, cleanName)

	compileToList(pt.outputTree.Root, pt.nodes)
}

func compileToList(arr *parse.ListNode, nodes []protoNode) {
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
		}
	}

}
