package bham

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

func newpipeline(s string) *parse.PipeNode {
	// take the simplest way of getting text/template to parse it
	// and then steal the result
	t, _ := template.New("mule").Parse("{{" + s + "}}")
	main := t.Tree.Root.Nodes[0]
	if an, ok := main.(*parse.ActionNode); ok {
		return an.Pipe
	} else {
		return nil
	}
}

func safeAction(s string) (*parse.ActionNode, error) {
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

func addEmbeddable(tn *parse.TextNode) []parse.Node {
	if strings.Contains(string(tn.Text), RightDelim) && strings.Contains(string(tn.Text), LeftDelim) {
		return []parse.Node{tn}
	} else {
		return []parse.Node{tn}
	}
}
