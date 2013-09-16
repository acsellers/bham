package bham

import (
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

func newMaybeTextNode(text string) []parse.Node {
	if strings.Contains(text, LeftDelim) && strings.Contains(text, RightDelim) {
		output := make([]parse.Node, 0)
		workingText := text
		for containsDelimeters(workingText) {
			index := strings.Index(workingText, LeftDelim)
			output = append(output, newTextNode(workingText[:index]))
			workingText = workingText[index:]

			index = strings.Index(workingText, RightDelim)
			pipeText := workingText[:index+len(RightDelim)]
			workingText = workingText[index+len(RightDelim):]

			action, e := safeAction(pipeText)
			if e != nil {
				output = append(output, newTextNode(pipeText+" "))
			} else {
				output = append(output, action)
			}
		}

		if workingText != "" {
			return append(output, newTextNode(workingText+" "))
		}
		return output
	} else {
		return []parse.Node{
			newTextNode(text + " "),
		}
	}
}

func newValueNode(val string) parse.Node {
	switch val {
	case "true":
		return &parse.BoolNode{
			NodeType: parse.NodeBool,
			True:     true,
		}
	case "false":
		return &parse.BoolNode{
			NodeType: parse.NodeBool,
			True:     false,
		}
	case "nil":
		return &parse.NilNode{}
	}
	panic("Can only call value node for true, false, nil")
}

func newFieldNode(field string) parse.Node {

	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						newBareFieldNode(field),
					},
				},
			},
		},
	}
}
func newBareFieldNode(field string) *parse.FieldNode {
	if field[0] == '.' {
		field = field[1:]
	}

	return &parse.FieldNode{
		NodeType: parse.NodeField,
		Ident:    strings.Split(field, "."),
	}
}

func newFunctionNode(command string) parse.Node {
	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args:     newBareFunctionNode(command),
				},
			},
		},
	}
}
func newBareFunctionNode(command string) []parse.Node {
	args := strings.Split(command, " ")
	nodeArgs := []parse.Node{}
	for _, arg := range args {
		switch arg[:1] {
		case ".":
			nodeArgs = append(nodeArgs,
				&parse.FieldNode{
					NodeType: parse.NodeField,
					Ident:    strings.Split(arg[1:], "."),
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

	return nodeArgs
}

func newStandaloneTag(content string) ([]parse.Node, error) {
	td, c, err := parseTag(content)
	if err != nil {
		return []parse.Node{}, err
	}
	return td.Nodes(c)
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
	var s runeScanner
	var current, attr, value, extra string
	td := newTagDescription()
	state := stateTag
	preamble, attributes := true, false

	s.Init(content)
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
				current = ""
				preamble = false
				attributes = true
			case ' ':
				td.Add(current, state)
				current = ""
				preamble = false
				attributes = false
				s.Reverse()
			case '=':
				td.Add(current, state)
				current = ""
				preamble = false
				attributes = false
				s.Reverse()
			default:
				current = current + s.TokenText()
			}
			tok = s.Scan()
		} else {
			if !attributes {
				for tok != scanner.EOF {
					extra = extra + s.TokenText()
					tok = s.Scan()
				}
				return td, extra, nil
			} else {
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
	}
	if current != "" {
		td.Add(current, state)
	}
	return td, extra, nil
}

func newTagDescription() tagDescription {
	return tagDescription{
		tag:        "div",
		attributes: make(map[string]string),
	}
}

type tagDescription struct {
	tag        string
	classes    []string
	idParts    []string
	attributes map[string]string
}

func (td *tagDescription) Add(content string, state int) {
	if content == "" {
		return
	}
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
	output := fmt.Sprintf("<%s", td.tag)
	if len(td.classes) > 0 {
		output = output + " class=\"" + strings.Join(td.classes, " ") + "\""
	}
	if len(td.idParts) > 0 {
		output = output + " id=\"" + strings.Join(td.idParts, IdJoin) + "\""
	}
	return output + ">"
}

func (td tagDescription) Close() string {
	return fmt.Sprintf("</%s>", td.tag)
}

func (td tagDescription) Nodes(content string) ([]parse.Node, error) {
	if content != "" {
		if content[0] == '=' {
			content = strings.TrimSpace(content[1:])
			node, err := parseTemplateCode(content)
			return []parse.Node{
				newTextNode(td.Opening()),
				&parse.ActionNode{
					NodeType: parse.NodeAction,
					Pipe:     node,
				},
				newTextNode(td.Close()),
			}, err
		} else {
			output := []parse.Node{
				newTextNode(td.Opening()),
			}
			output = append(output, newMaybeTextNode(content)...)
			return append(output, newTextNode(td.Close())), nil
		}
	} else {
		return []parse.Node{
			newTextNode(td.Opening()),
			newTextNode(td.Close()),
		}, nil
	}
}
