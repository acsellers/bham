package bham

var Filters = []FilterHandler{
	FilterHandler{
		Trigger: ":javascript",
		Open:    `<script type="text/javascript">`,
		Close:   "</script>",
		Handler: func(s string) { return s },
	},
}

type FilterHandler struct {
	Trigger     string
	Open, Close string
	Handler     Transformer
}

type Transformer func(string) string

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
