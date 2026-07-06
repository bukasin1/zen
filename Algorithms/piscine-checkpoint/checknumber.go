package main

func CheckNumber(arg string) bool {
	for _, ch := range arg {
		if '0' <= ch && ch <= '9' {
			return true
		}
	}
	return false
}
