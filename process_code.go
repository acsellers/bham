package bham

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template/parse"
	"unicode"
)

func processCode(s string) (*parse.ActionNode, error) {
	var chars []rune
	for _, char := range s {
		chars = append(chars, char)
	}

	current, last := 0, len(chars)-1
	continuing := true
	work := new(bytes.Buffer)
	var funcName string
	var args []string
	var nodeArgs []parse.Node
	var nodes []*parse.CommandNode

	{
	begin_command:
		for current <= last && unicode.IsSpace(chars[current]) {
			current++
		}
		for current <= last && !unicode.IsSpace(chars[current]) {
			work.WriteRune(chars[current])
			current++
		}
		funcName = work.String()
		work.Reset()
		for current <= last && unicode.IsSpace(chars[current]) {
			current++
		}
		if current >= last {
			goto complete
		}
		goto choose_after

	func_arg:
		if chars[current] == '"' {
			work.WriteRune(chars[current])
			current++
			for current <= last && chars[current] != '"' && chars[current-1] != '\\' {
				work.WriteRune(chars[current])
				current++
			}
			if chars[current] != '"' {
				return nil, fmt.Errorf("Unterminated string: %s", work.String())
			}
			work.WriteRune(chars[current])
			current++
		} else {
			for current <= last && !unicode.IsSpace(chars[current]) {
				work.WriteRune(chars[current])
				current++
			}
		}

		for current <= last && unicode.IsSpace(chars[current]) {
			current++
		}
		args = append(args, work.String())
		work.Reset()

		if current >= last {
			goto complete
		}

	choose_after:
		if chars[current] == '|' {
			goto push_func
		} else {
			goto func_arg
		}

	complete:
		continuing = false

	push_func:
		nodeArgs = []parse.Node{}
		nodeArgs = append(nodeArgs,
			&parse.IdentifierNode{
				NodeType: parse.NodeIdentifier,
				Ident:    funcName,
			})
		for _, arg := range args {
			nodeArgs = append(nodeArgs, parseFuncArg(arg))
		}
		funcName = ""
		nodes = append(nodes, &parse.CommandNode{
			NodeType: parse.NodeCommand,
			Args:     nodeArgs,
		})

		args = []string{}

		if continuing {
			goto begin_command
		}
	}

	if current < last {
		return nil, fmt.Errorf("Could not complete parse")
	}
	if funcName != "" {
		return nil, fmt.Errorf("Parse was not able to complete")
	}

	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds:     nodes,
		},
	}, nil
}

func parseFuncArg(arg string) parse.Node {
	// fields (.first.second)
	// variables ($first.second)
	// strings ("value")
	switch arg[0] {
	case '.':
		if arg == "." {
			return &parse.DotNode{}
		} else {
			return &parse.FieldNode{
				NodeType: parse.NodeField,
				Ident:    strings.Split(arg[1:], "."),
			}
		}
	case '$':
		return &parse.VariableNode{
			NodeType: parse.NodeVariable,
			Ident:    strings.Split(arg, "."),
		}
	case '"':
		return &parse.StringNode{
			NodeType: parse.NodeString,
			Quoted:   arg,
			Text:     arg[1 : len(arg)-1],
		}
	}

	// bool variables (true || false)
	// the value nil
	switch arg {
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

	// numeric
	if node, ok := parseNumeric(arg); ok {
		return node

		// function names (blah, url, etc)
	} else {
		return &parse.IdentifierNode{
			NodeType: parse.NodeIdentifier,
			Ident:    arg,
		}
	}
}

func parseNumeric(num string) (*parse.NumberNode, bool) {
	node := &parse.NumberNode{
		NodeType: parse.NodeNumber,
		Text:     num,
	}

	// Following code is adapted from text/template's newNumber function
	// Handle all int's first
	if num[0] != '-' {
		u, err := strconv.ParseUint(num, 0, 64) // will fail for -0; fixed below.
		if err == nil {
			node.IsUint = true
			node.Uint64 = u
			node.IsFloat = true
			node.Float64 = float64(u)
		}
	}

	i, err := strconv.ParseInt(num, 0, 64)
	if err == nil {
		node.IsInt = true
		node.Int64 = i
		if i == 0 {
			node.IsUint = true // in case of -0.
			node.Uint64 = 0
		}
		node.IsFloat = true
		node.Float64 = float64(i)
		return node, true
	}
	if node.IsUint {
		return node, true
	}

	// handle all floats as floats only
	// text template will allow you to turn integer floats into
	// ints/uints, I'm just not feeling it
	f, err := strconv.ParseFloat(num, 64)
	if err == nil {
		node.IsFloat = true
		node.Float64 = f
		return node, true
	}

	return nil, false
}