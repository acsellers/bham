package bham

import (
	"bytes"
	"fmt"
	"strings"
	"text/scanner"
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
	if field[0] == '.' {
		field = field[1:]
	}

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

func newFunctionNode(command string) parse.Node {
	args := strings.Split(command, " ")
	nodeArgs := []parse.Node{}
	for _, arg := range args {
		switch arg[:1] {
		case ".":
			nodeArgs = append(nodeArgs,
				&parse.FieldNode{
					NodeType: parse.NodeField,
					Ident:    strings.Split(arg, "."),
				},
			)
		case "$":
			nodeArgs = append(nodeArgs,
				&parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident:    strings.Split(arg, "."),
				},
			)
		default:
			nodeArgs = append(nodeArgs,
				&parse.IdentifierNode{
					NodeType: parse.NodeIdentifier,
					Ident:    arg,
				},
			)
		}
	}

	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args:     nodeArgs,
				},
			},
		},
	}
}

func newStandaloneTag(content string) ([]parse.Node, error) {
	_, _, err := parseTag(content)
	if err != nil {
		return []parse.Node{}, err
	}
	return []parse.Node{}, nil
}

const (
	stateNull = iota
	stateTag
	stateId
	stateClass
	stateAttr
	stateValue
	stateBridge
	stateSpace
	stateDone
)

func parseTag(content string) (tagDescription, string, error) {
	var s scanner.Scanner
	var current, attr, value, extra string
	td := newTagDescription()
	state := stateTag
	preamble, attributes := true, false

	s.Init(bytes.NewBufferString(content))
	tok := s.Scan()
	for tok != scanner.EOF {
		if preamble {
			switch tok {
			case '.':
				td.Add(current, state)
				state = stateClass
				current = ""
			case '#':
				td.Add(current, state)
				state = stateId
				current = ""
			case '%':
				if state != stateTag || current != "" {
					return td, "", fmt.Errorf("Tags must start the element declaration")
				}
			case '(':
				td.Add(current, state)
				preamble = false
				attributes = true
			case ' ':
				preamble = false
				attributes = false
			default:
				current = current + s.TokenText()
			}
			tok = s.Scan()
		} else {
			if !attributes {
				tok = s.Scan()
				for tok != scanner.EOF {
					extra = extra + s.TokenText()
					tok = s.Scan()
				}
				return td, extra, nil
			}
			switch tok {
			case '=':
				if state == stateValue {
					value = value + "="
				} else {
					state = stateBridge
				}
			case '"':
				switch s.Peek() {
				case ')':
					state = stateDone
					s.Scan()
				case ' ':
					state = stateSpace
				}
			case ' ':
				if state == stateValue {
					value = value + " "
				}
			case ',':
				switch state {
				case stateValue:
					value = value + ","
				case stateAttr:
					td.attributes[attr] = ""
					state = stateSpace
				}
			case ')':
				switch state {
				case stateDone:
					tok = s.Scan()
					for tok != scanner.EOF {
						extra = extra + s.TokenText()
						tok = s.Scan()
					}
					return td, extra, nil
				case stateValue:
					value = value + s.TokenText()
				}
			default:
				switch state {
				case stateAttr:
					switch s.Peek() {
					case ' ':
						state = stateSpace
						td.attributes[attr] = ""
					case ')':
						state = stateDone
						for tok != scanner.EOF {
							extra = extra + s.TokenText()
							tok = s.Scan()
						}
						return td, extra, nil
					default:
						attr = attr + s.TokenText()
					}
				case stateValue:
					value = value + s.TokenText()
				case stateSpace:
					state = stateAttr
				case stateBridge:
					value = value + s.TokenText()
				}
			}
		}
	}
	return td, extra, nil
}

func newTagDescription() tagDescription {
	return tagDescription{
		attributes: make(map[string]string),
	}
}

type tagDescription struct {
	tag        string
	classes    []string
	idParts    []string
	attributes map[string]string
}

func (td tagDescription) Add(content string, state int) {
	switch state {
	case stateClass:
		td.classes = append(td.classes, content)
	case stateId:
		td.idParts = append(td.idParts, content)
	case stateTag:
		td.tag = content
	}
}
func (td tagDescription) Opening() string {
	return ""
}

func (td tagDescription) Close() string {
	return ""
}