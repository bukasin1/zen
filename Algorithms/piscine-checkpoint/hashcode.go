package main

func HashCode(dec string) string {
	var codedstr string

	for _, ch := range dec {
		code := (int(ch) + len(dec)) % 127
		if code <= 0 {
			code += 33
		}

		codedstr += string(rune(code))
	}

	return codedstr
}
