package bham

import (
	"fmt"
	"strings"
	"text/template/parse"
)

func (pt *protoTree) newListNode(listarea []token) *parse.ListNode {
	listNode := new(parse.ListNode)

	var currentIndex, localIndex int
	var currentToken token

	for currentIndex < len(listarea) {
		currentToken = listarea[currentIndex]
		switch currentToken.purpose {
		case pse_text, pse_tag:
			textNode := new(parse.TextNode)
			textNode.NodeType = parse.NodeText

			localIndex = currentIndex
			for localIndex < len(listarea) && listarea[localIndex].textual() {
				localIndex++
			}
			texts := []string{""}
			lastPurpose := pse_tag
			for _, token := range listarea[currentIndex:localIndex] {
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
			listNode.Nodes = append(listNode.Nodes, addEmbeddable(textNode)...)
			currentIndex = localIndex
		case pse_if:
			ifNode := &parse.IfNode{
				newBranchNode(parse.NodeIf, currentToken.content),
			}
			listNode.Nodes = append(listNode.Nodes, ifNode)

			localIndex = currentIndex + 1
			for listarea[localIndex].parent() != currentIndex {
				localIndex++
			}
			ifNode.BranchNode.List = pt.newListNode(
				listarea[currentIndex+1 : localIndex],
			)

			if listarea[localIndex].purpose == pse_else {
				currentIndex = localIndex
				localIndex = currentIndex + 1
				for listarea[localIndex].parent() != currentIndex {
					localIndex++
				}
				ifNode.BranchNode.ElseList = pt.newListNode(
					listarea[currentIndex+1 : localIndex],
				)
			}

			currentIndex = localIndex + 1
		case pse_range:
			rangeNode := &parse.RangeNode{
				newBranchNode(parse.NodeRange, currentToken.content),
			}
			listNode.Nodes = append(listNode.Nodes, rangeNode)

			localIndex = currentIndex + 1
			for listarea[localIndex].parent() != currentIndex {
				localIndex++
			}
			rangeNode.BranchNode.List = pt.newListNode(
				listarea[currentIndex+1 : localIndex],
			)

			if listarea[localIndex].purpose == pse_else {
				currentIndex = localIndex
				localIndex = currentIndex + 1
				for listarea[localIndex].parent() != currentIndex {
					localIndex++
				}
				rangeNode.BranchNode.ElseList = pt.newListNode(
					listarea[currentIndex+1 : localIndex],
				)
			}

			currentIndex = localIndex + 1
		case pse_with:
			withNode := &parse.WithNode{
				newBranchNode(parse.NodeWith, currentToken.content),
			}

			listNode.Nodes = append(listNode.Nodes, withNode)

			localIndex = currentIndex + 1
			for listarea[localIndex].parent() != currentIndex {
				localIndex++
			}
			withNode.BranchNode.List = pt.newListNode(
				listarea[currentIndex+1 : localIndex],
			)

			currentIndex = localIndex + 1
		case pse_exe:
			templ := listarea[currentIndex].content
			an, e := safeAction("{{" + templ + "}}")
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
