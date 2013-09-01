package bham

import (
	"bufio"
	"bytes"
	"strings"
)

func (pt *protoTree) lex() {
	scanner := bufio.NewScanner(bytes.NewBufferString(pt.source))
	var line, content string
	var currentLevel int
	var currentLine int
	for scanner.Scan() {
		currentLine++
		line = scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		level, content = level(line)
		if currentLevel+1 >= level {
			pt.lineList = append(pt.lineList, templateLine{level, content})
		} else {
			pt.err = fmt.Errorf("Line %d is overindented", currentLine)
		}
	}
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
