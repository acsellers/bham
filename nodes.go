package bham

const (
	identRaw = iota
	identText
)

func (pt *protoTree) insertRaw(content string, level int) {
	pt.nodes = append(pt.nodes, protoNode{
		level:      level,
		identifier: identRaw,
		content:    content,
	})
}

func (pt *protoTree) insertText(line templateLine) {
	pt.nodes = append(pt.nodes, protoNode{
		level:      line.indentation,
		identifier: identText,
		content:    line.content,
	})
}
