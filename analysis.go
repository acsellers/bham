package bham

func (pt *protoTree) analyze() {
	previousline := templateLine{-1, ""}
	currentIndex := 0
	for currentIndex < len(pt.lineList) {
		if pt.err != nil {
			return
		}

		line := pt.lineList[currentIndex]
		if line.indentation > previousline.indentation+1 {
			pt.err = fmt.Errorf("Line %d is indented more than necessary", currentIndex+1)
			return
		}

		switch {
		case line.accept("%.#"):
			currentIndex = pt.tagLike(currentIndex)
			continue
		case line.accept("=-"):
			currentIndex = pt.actionableLine(currentIndex)
			continue
		case line.accept(":"):
			for _, handler := range Filters {
				if line.prefix(handler.Trigger) {
					currentIndex = pt.followHandler(currentIndex+1, handler, line.indentation)
					continue
				}
			}
			pt.err = fmt.Errorf("Bad handler: %s", line.content)
			return
		case line.prefix("!!!"):
			pt.insertDoctype(line)
		default:
			pt.insertText(line)
		}

		currentIndex++
	}
}

func (pt *protoTree) insertDoctype(line templateLine) {
	doctype, ok := Doctypes[strings.TrimSpace(line.content[3:])]
	if ok {
		pt.insertRaw(doctype, line.indentation)
	} else {
		pt.err = fmt.Errorf("Bad doctype, details: '%s'", line.content)
	}
}

func (pt *protoTree) followHandler(startIndex int, handler FilterHandler, level int) int {
	lines := make([]string, 0)
	index := startIndex
	base := pt.lineList[startIndex].indentation
	for base <= pt.lineList[index].indentation {
		diff := base - pt.lineList[index].indentation
		lines = append(lines, pad(diff)+pt.lineList[index].content)
		index++
	}
	pt.insertRaw(handler.Open + handler.Handler(strings.Join(lines, "\n")) + handler.Close)

	return index
}

func pad(indent int) string {
	var output string
	for i := 0; i < indent; i++ {
		output = output + "  "
	}
	return output
}

func (pt *protoTree) actionableLine(currentIndex int) int {
	line := pt.lineList[currentIndex]
	text := line.content

	return currentIndex + 1
}

func (pt *protoTree) tagLike(currentIndex int) int {
}
