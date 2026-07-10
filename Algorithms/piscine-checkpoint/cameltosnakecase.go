package main

func isCapital(ch rune) bool {
	return 'A' <= ch && 'Z' >= ch
}

func CamelToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	lastChar := s[len(s)-1]
	if isCapital(rune(lastChar)) {
		return s
	}
	var snake_case string

	for i, ch := range s {
		if ch >= '0' && ch <= '9' {
			return s
		}
		if isCapital(ch) {
			if i+1 < len(s) && isCapital(rune(s[i+1])) {
				return s
			}
			if i > 0 {
				// snake_case += "_" + string(ch+'a'-'A')
				snake_case += "_" + string(ch)
			} else {
				// snake_case += string(ch + 'a' - 'A')
				snake_case += string(ch)
			}
		} else {
			snake_case += string(ch)
		}
	}

	return snake_case
}
