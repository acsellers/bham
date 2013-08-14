package bham

import (
	"strings"
	"text/template/parse"
)

// parse will return a parse tree containing a single
func Parse(name, text string) (map[string]*parse.Tree, error) {
	proto := &protoTree{source: text}
	proto.tokenize()
	i := strings.Index(name, ".bham")

	return map[string]*parse.Tree{
		name[:i] + name[i+5:]: proto.tree(),
	}, proto.err
}

type protoTree struct {
	name      string
	source    string
	tokenList []token
	err       error
}

func (pt *protoTree) tree() *parse.Tree {
	if pt.err != nil {
		return nil
	}
	tree := &parse.Tree{Root: pt.newListNode(pt.tokenList)}
	return tree
}
