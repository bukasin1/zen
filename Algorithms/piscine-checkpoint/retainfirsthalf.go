package main

func RetainFirstHalf(str string) string {
	if len(str) <= 1 {
		return str
	}

	halfLen := len(str) / 2

	return str[:halfLen]
}
