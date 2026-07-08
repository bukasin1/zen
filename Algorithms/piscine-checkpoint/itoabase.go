package main

func ItoaBase(value, base int) string {
	if value == 0 {
		return "0"
	}

	var isnegative bool
	if value < 0 {
		value = -value
		isnegative = true
	}

	var valStr string
	for value > 0 {
		rem := value % base
		var remstr rune
		if rem > 9 {
			remstr = rune('A' + rem - 10)
		} else {
			remstr = rune('0' + rem)
		}
		valStr = string(remstr) + valStr
		value /= base
	}

	if isnegative {
		valStr = "-" + valStr
	}

	return valStr
}
