package bham

import (
	"text/scanner"
)

type runeScanner struct {
	chars   []rune
	current int
}

func (rs *runeScanner) Init(s string) {
	for _, char := range s {
		rs.chars = append(rs.chars, char)
	}
}

func (rs *runeScanner) Scan() rune {
	if rs.current < len(rs.chars) {
		rs.current++
		return rs.chars[rs.current-1]
	}

	return scanner.EOF
}

func (rs runeScanner) TokenText() string {
	return string([]rune{rs.chars[rs.current-1]})
}

func (rs runeScanner) Peek() rune {
	if rs.current < len(rs.chars) {
		return rs.chars[rs.current]
	}
	return scanner.EOF
}

func (rs *runeScanner) Reverse() {
	rs.current--
}
