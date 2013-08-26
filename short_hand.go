package bham

var ShortHands = map[string]ShortHandCode{
	"javascript": ShortHandCode{
		Open:  `<script type="text/javascript">`,
		Close: "</script>",
	},
}

type ShortHandCode struct {
	Open, Close string
}

func shortHandOpen(sh string) token {
	output := token{
		purpose: pse_tag,
	}
	if handler, ok := ShortHands[sh]; ok {
		output.content = handler.Open
	} else {
		output.content = "<" + sh + ">"
	}
	return output
}

func shortHandClose(sh string) token {
	output := token{
		purpose: pse_tag,
	}
	if handler, ok := ShortHands[sh]; ok {
		output.content = handler.Close
	} else {
		output.content = "</" + sh + ">"
	}
	return output
}
