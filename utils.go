package bham

func findAttrs(s string) (string, string) {
	var openings int
	for i, r := range s {
		switch r {
		case '(':
			openings++
		case ')':
			openings--
			if openings == 0 {
				return s[1:i], s[i+1:]
			}
		}
	}
	return "", s
}
