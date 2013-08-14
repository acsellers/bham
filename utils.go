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
