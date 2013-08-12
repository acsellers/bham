package bham

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const (
	pse_text = iota
	pse_tag
	pse_if
	pse_else
	pse_end
	pse_range
	pse_with
	pse_decl
	pse_exe
)

func (pt *protoTree) tokenize() error {
	posts := make([]token, 0, 64)
	scanner := bufio.NewScanner(bytes.NewBufferString(pt.source))
	var text, currentTag string
	var currentLevel, currentLine, lineLevel int
	for scanner.Scan() {
		currentLine++

		text = scanner.Text()
		if text == "" || strings.TrimSpace(text) == "" {
			continue
		}

		lineLevel, text = level(text)
		for currentLevel >= lineLevel && currentLevel > 0 {
			pt.tokenList = append(
				pt.tokenList,
				posts[len(posts)-1],
			)
			posts = posts[:len(posts)-1]
			currentLevel--
		}
		if lineLevel-1 > currentLevel {
			return fmt.Errorf("Line %d is indented more than necessary (%d) from the previous line %d", currentLine, lineLevel, currentLevel)
		}

		if tag.MatchString(text) {
			currentTag = tag.FindStringSubmatch(text)[1]
			pt.tokenList = append(
				pt.tokenList,
				token{
					content: "<" + currentTag + ">",
					purpose: pse_tag,
				},
			)

			posts = append(posts, token{
				content: "</" + currentTag + ">",
				purpose: pse_tag,
			})
			text = text[len(currentTag)+1:]
		} else {
			if strings.HasPrefix(text, LineDelim) {
				trimText := strings.TrimSpace(text[len(LineDelim):])
				switch {
				case strings.HasPrefix(trimText, "if "):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: strings.TrimPrefix(trimText, "if "),
							purpose: pse_if,
						},
					)
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case trimText == "else":
					pt.tokenList[len(pt.tokenList)-1].purpose = pse_else
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case strings.HasPrefix(trimText, "range "):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: strings.TrimPrefix(trimText, "range "),
							purpose: pse_range,
						},
					)
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case strings.HasPrefix(trimText, "with "):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: strings.TrimPrefix(trimText, "with "),
							purpose: pse_with,
						},
					)
					posts = append(posts, token{purpose: pse_end, previous: len(pt.tokenList) - 1})
				case varDecl.MatchString(trimText):
					pt.tokenList = append(
						pt.tokenList,
						token{
							content: trimText,
							purpose: pse_decl,
						},
					)
				default:
					pt.tokenList = append(pt.tokenList, token{content: trimText, purpose: pse_exe})
				}
			} else {
				pt.tokenList = append(pt.tokenList, token{content: text})
			}
			text = ""
		}
		if text == "" {
			currentLevel = lineLevel
			continue
		}
		if text[0] == ' ' {
			pt.tokenList = append(pt.tokenList, token{content: text[1:]})
			currentLevel = lineLevel
			continue
		}
	}
	for i, _ := range posts {
		pt.tokenList = append(pt.tokenList, posts[len(posts)-i-1])
	}

	return nil
}

func (pt *protoTree) compact() {
	newTokenList := make([]token, 0, len(pt.tokenList))
	var appendToken token
	appendToken.purpose = pse_text
	for _, currentToken := range pt.tokenList {
		if currentToken.purpose == pse_text {
			appendToken.content += currentToken.content
		} else {
			newTokenList = append(newTokenList,
				appendToken,
				currentToken,
			)
			appendToken.content = ""
		}
	}
	pt.tokenList = newTokenList
}

type token struct {
	content  string
	purpose  int
	previous int
}

func (t token) parent() int {
	switch t.purpose {
	case pse_end:
		return t.previous
	case pse_else:
		return t.previous
	default:
		return -1
	}
}

func (t token) textual() bool {
	return t.purpose == pse_text || t.purpose == pse_tag
}
