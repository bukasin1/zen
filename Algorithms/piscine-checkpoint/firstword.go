package main

func FirstWord(s string) string {
	var word string
	leadingSpace := true

	for _, ch := range s {
		if ch == ' ' {
			if !leadingSpace {
				return word + "\n"
			} else {
				continue
			}
		} else {
			word += string(ch)
			leadingSpace = false
		}
	}

	return word + "\n"
}
