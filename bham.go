package bham

import (
	"strings"
	"text/template/parse"
)

var Doctypes = map[string]string{
	"":             `<!DOCTYPE html>`,
	"Transitional": `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`,
	"Strict":       `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
	"Frameset":     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">`,
	"5":            `<!DOCTYPE html>`,
	"1.1":          `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
	"Basic":        `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN" "http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd">`,
	"Mobile":       `<!DOCTYPE html PUBLIC "-//WAPFORUM//DTD XHTML Mobile 1.2//EN" "http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd">`,
	"RDFa":         `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML+RDFa 1.0//EN" "http://www.w3.org/MarkUp/DTD/xhtml-rdfa-1.dtd">`,
}

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
	lineList  []templateLine
	nodes     []protoNode
	tokenList []token
	err       error
}

type templateLine struct {
	indentation int
	content     string
}

func (t templateLine) accept(chars string) bool {
	for _, s := range chars {
		if s == t.content[0] {
			return true
		}
	}
	return false
}
func (t templateLine) prefix(str string) bool {
	return len(t.content) >= len(str) && t.content[:len(str)] == str
}

type protoNode struct {
	level      int
	identifier int
	content    string
}

func (pt *protoTree) tree() *parse.Tree {
	if pt.err != nil {
		return nil
	}
	tree := &parse.Tree{Root: pt.newListNode(pt.tokenList)}
	return tree
}
