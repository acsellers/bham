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

	// To add multiple id declarations, the outputter puts them together
	// with a join string, by default this is an underscore
	IdJoin = "_"

	// Like the template library, you need to be able to set code delimeters
	LeftDelim  = "{{"
	RightDelim = "}}"
	LineDelim  = "="
)
var (
	tag      = regexp.MustCompile("^%([a-zA-Z0-9]+)")
	varDecl  = regexp.MustCompile("^\\$[a-zA-Z0-9]+ :=")
	idClass  = regexp.MustCompile("^([\\.#][a-zA-Z0-9-_]+)")
	idClass2 = regexp.MustCompile("([\\.#][a-zA-Z0-9-_]+)")
	varUse   = regexp.MustCompile("(\\$[a-zA-Z0-9]+)")
)

// parse will return a parse tree containing a single
func Parse(name, text string) (map[string]*parse.Tree, error) {
	proto := &protoTree{source: text}
	proto.tokenize()
	i := strings.Index(name, ".bham")

	return map[string]*parse.Tree{name[:i] + name[i+5:]: proto.treeify()}, proto.err
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

	var currentIndex, textIndex, ifIndex, rangeIndex, withIndex int
	var currentToken token

	for currentIndex < len(listarea) {
		currentToken = listarea[currentIndex]
		switch currentToken.purpose {
		case pse_text, pse_tag:
			textNode := new(parse.TextNode)
			textNode.NodeType = parse.NodeText
			listNode.Nodes = append(listNode.Nodes, textNode)

			textIndex = currentIndex
			for textIndex < len(listarea) && listarea[textIndex].textual() {
				textIndex++
			}
			texts := []string{""}
			lastPurpose := pse_tag
			for _, token := range listarea[currentIndex:textIndex] {
				if lastPurpose == pse_tag {
					texts[len(texts)-1] += token.strcontent()
				} else {
					if token.purpose == pse_tag {
						texts[len(texts)-1] += token.strcontent()
					} else {
						texts = append(texts, token.strcontent())
					}
				}
				lastPurpose = token.purpose
			}
			textNode.Text = append(textNode.Text, []byte(strings.Join(texts, " "))...)
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
			rangeNode := &parse.RangeNode{
				parse.BranchNode{
					NodeType: parse.NodeRange,
					Pipe:     pt.pipeify(currentToken.content),
				},
			}
			listNode.Nodes = append(listNode.Nodes, rangeNode)

			rangeIndex = currentIndex + 1
			for listarea[rangeIndex].parent() != currentIndex {
				rangeIndex++
			}
			rangeNode.BranchNode.List = pt.listify(
				listarea[currentIndex+1 : rangeIndex],
			)

			if listarea[rangeIndex].purpose == pse_else {
				currentIndex = rangeIndex
				rangeIndex = currentIndex + 1
				for listarea[rangeIndex].parent() != currentIndex {
					rangeIndex++
				}
				rangeNode.BranchNode.ElseList = pt.listify(
					listarea[currentIndex+1 : rangeIndex],
				)
			}

			currentIndex = rangeIndex + 1
		case pse_with:
			withNode := &parse.WithNode{
				parse.BranchNode{
					NodeType: parse.NodeWith,
					Pipe:     pt.pipeify(currentToken.content),
				},
			}
			listNode.Nodes = append(listNode.Nodes, withNode)

			withIndex = currentIndex + 1
			for listarea[withIndex].parent() != currentIndex {
				withIndex++
			}
			withNode.BranchNode.List = pt.listify(
				listarea[currentIndex+1 : withIndex],
			)

			currentIndex = withIndex + 1
		case pse_exe:
			templ := listarea[currentIndex].content
			an, e := pt.safeActionify("{{" + templ + "}}")
			if e != nil {
				listarea[currentIndex].purpose = pse_text
				listarea[currentIndex].content = "= " + listarea[currentIndex].content
				fmt.Printf("Couldn't parse %s, got %v\n", templ, e)
				continue
			}
			listNode.Nodes = append(listNode.Nodes, an)
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

func (pt *protoTree) safeActionify(s string) (*parse.ActionNode, error) {
	// take the simplest way of getting text/template to parse it
	// and then steal the result
	if varUse.MatchString(s) {
		for _, varUser := range varUse.FindAllStringSubmatch(s, -1) {
			s = "{{ " + varUser[0] + " := 0 }}" + s
		}
	}
	t, e := template.New("mule").Parse(s)
	if e != nil {
		return nil, e
	}
	main := t.Tree.Root.Nodes[len(t.Tree.Root.Nodes)-1]
	if an, ok := main.(*parse.ActionNode); ok {
		return an, nil
	} else {
		return nil, fmt.Errorf("Couldn't find action node")
	}
}
