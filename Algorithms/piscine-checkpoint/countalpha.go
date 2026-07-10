package main

func CountAlpha(s string) int {
	alphaCount := 0
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			alphaCount++
		}
	}

	return alphaCount
}
