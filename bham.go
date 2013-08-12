package bham

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"text/template/parse"
)

var (
	// Strict determines whether only tabs will be considered
	// as indentation operators (Strict == true) or whether
	// two spaces can be counted as an indentation operator
	// (Strict == false), this is included for haml
	// semi-comapibility
	Strict bool

	// Like the template library, you need to be able to set code delimeters
	LeftDelim  = "{{"
	RightDelim = "}}"
	LineDelim  = "="

	// Since bham will likely need to break templates around a yield call
	// you may set suffixes for the file name to be set in the
	// return map. By default a file named "test.bham" with a yield
	// command would have keys for both test-upper.bham and test-lower.bham
	YieldFirstSuffix  = "-upper"
	YieldSecondSuffix = "-lower"
)
var tag = regexp.MustCompile("^%([a-zA-Z0-9]+)")

// parse will return a parse tree containing a single
func Parse(name, text string) (map[string]*parse.Tree, error) {
	proto := &protoTree{source: text}
	proto.tokenize()
	i := strings.Index(name, ".bham")

	return map[string]*parse.Tree{name[:i] + name[i+5:]: proto.treeify()}, proto.err
}

func ParseFiles(filenames ...string) (*parse.Tree, error) {
	return nil, nil
}

func ParseGlob(pattern string) (*parse.Tree, error) {
	return nil, nil
}

type protoTree struct {
	name      string
	source    string
	tokenList []token
	err       error
}

func level(s string) (int, string) {
	var currentLevel int
	for {
		switch s[0] {
		case ' ':
			if s[1] == ' ' {
				s = s[2:]
			} else {
				return currentLevel, s
			}
		case '\t':
			s = s[1:]
		default:
			return currentLevel, s
		}
		currentLevel++
	}
}

func (pt *protoTree) treeify() *parse.Tree {
	if pt.err != nil {
		return nil
	}
	tree := &parse.Tree{Root: pt.listify(pt.tokenList)}
	return tree
}

func (pt *protoTree) listify(listarea []token) *parse.ListNode {
	listNode := new(parse.ListNode)

	var currentIndex, textIndex, ifIndex int
	var currentToken token

	for currentIndex < len(listarea) {
		currentToken = listarea[currentIndex]
		switch currentToken.purpose {
		case pse_text:
			textNode := new(parse.TextNode)
			textNode.NodeType = parse.NodeText
			listNode.Nodes = append(listNode.Nodes, textNode)

			textIndex = currentIndex
			for textIndex < len(listarea) && listarea[textIndex].purpose == pse_text {
				textIndex++
			}
			for _, token := range listarea[currentIndex:textIndex] {
				textNode.Text = append(textNode.Text, []byte(token.content)...)
			}
			currentIndex = textIndex
		case pse_if:
			ifNode := &parse.IfNode{
				parse.BranchNode{
					NodeType: parse.NodeIf,
					Pipe:     pt.pipeify(currentToken.content),
				},
			}
			listNode.Nodes = append(listNode.Nodes, ifNode)

			ifIndex = currentIndex + 1
			for listarea[ifIndex].parent() != currentIndex {
				ifIndex++
			}
			ifNode.BranchNode.List = pt.listify(
				listarea[currentIndex+1 : ifIndex],
			)

			if listarea[ifIndex].purpose == pse_else {
				currentIndex = ifIndex
				ifIndex = currentIndex + 1
				for listarea[ifIndex].parent() != currentIndex {
					ifIndex++
				}
				ifNode.BranchNode.ElseList = pt.listify(
					listarea[currentIndex+1 : ifIndex],
				)
			}

			currentIndex = ifIndex + 1
		case pse_range:
			fmt.Println("ERROR: range not written yet")
			currentIndex++
		case pse_with:
			fmt.Println("ERROR: with not written yet")
			currentIndex++
		default:
			fmt.Println("ERROR: token not recognized", listarea[currentIndex])
			currentIndex++
		}
	}
	return listNode
}

func (pt *protoTree) pipeify(s string) *parse.PipeNode {
	// take the simplest way of getting text/template to parse it
	// and then steal the result
	t, _ := template.New("mule").Parse("{{" + s + "}}")
	main := t.Tree.Root.Nodes[0]
	if an, ok := main.(*parse.ActionNode); ok {
		return an.Pipe
	} else {
		panic("could not locate type")
	}
}
