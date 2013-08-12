package bham

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
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
	proto.classify()
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
	name       string
	source     string
	tokenList  []string
	classified []classifiers
	err        error
}

func (pt *protoTree) tokenize() error {
	posts := make([]string, 0, 64)
	scanner := bufio.NewScanner(bytes.NewBufferString(pt.source))
	var text, currentTag string
	var currentLevel, currentLine, lineLevel int
	for scanner.Scan() {
		currentLine++

		text = scanner.Text()
		if text == "" {
			continue
		}

		lineLevel, text = level(text)
		if text == "" {
			continue
		}

		for currentLevel >= lineLevel && currentLevel > 0 {
			pt.tokenList = append(pt.tokenList, posts[len(posts)-1])
			posts = posts[:len(posts)-1]
		}
		if lineLevel-1 > currentLevel {
			return fmt.Errorf("Line %d is indented more than necessary (%d) from the previous line %d", currentLine, lineLevel, currentLevel)
		}
		if tag.MatchString(text) {
			currentTag = tag.FindStringSubmatch(text)[1]
			pt.tokenList = append(pt.tokenList, "<"+currentTag+">")
			posts = append(posts, "</"+currentTag+">")
			text = text[len(currentTag)+1:]
		} else {
			pt.tokenList = append(pt.tokenList, text)
			text = ""
		}
		if text == "" {
			currentLevel = lineLevel
			continue
		}
		if text[0] == ' ' {
			pt.tokenList = append(pt.tokenList, text[1:])
			currentLevel = lineLevel
			continue
		}
	}
	for i, _ := range posts {
		pt.tokenList = append(pt.tokenList, posts[len(posts)-i-1])
	}

	return nil
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

func (pt *protoTree) classify() {
	if pt.err != nil {
		return
	}
	var currentBlock string
	for _, token := range pt.tokenList {
		switch {
		case strings.HasPrefix(token, LeftDelim) && strings.HasSuffix(token, RightDelim):
			pt.classified = append(pt.classified, textClassifier{currentBlock})
			currentBlock = ""
			pt.classified = append(pt.classified, fieldClassifier{token})
		default:
			currentBlock += token
		}
	}
	if currentBlock != "" {
		pt.classified = append(pt.classified, textClassifier{currentBlock})
	}
}

type classifiers interface {
	Executable() bool
	String() string
}

type textClassifier struct {
	data string
}

func (tc textClassifier) Executable() bool {
	return false
}

func (tc textClassifier) String() string {
	return tc.data
}

type fieldClassifier struct {
	data string
}

func (fc fieldClassifier) Executable() bool {
	return true
}
func (fc fieldClassifier) String() string {
	return fc.data
}

func (pt *protoTree) treeify() *parse.Tree {
	if pt.err != nil {
		return nil
	}
	tree := &parse.Tree{Root: &parse.ListNode{
		NodeType: parse.NodeList,
		Pos:      0,
	},
	}
	var currentPos int
	for _, classifier := range pt.classified {
		if classifier.Executable() {
		} else {
			tree.Root.Nodes = append(tree.Root.Nodes, &parse.TextNode{
				NodeType: parse.NodeText,
				Pos:      parse.Pos(currentPos),
				Text:     []byte(classifier.String()),
			})
		}
		currentPos += len(classifier.String())
	}

	return tree
}
