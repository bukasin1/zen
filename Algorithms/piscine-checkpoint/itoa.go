package main

// func itoa(n int) string {
// 	if n == 0 {

// 	}
// 	rem := n%10

// 	itoa(n/10)

// 	return ""
// }

func Itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var isnegative bool
	if n < 0 {
		n = -n
		isnegative = true
	}
	// fmt.Println("n to be processed:", n, isnegative)

	var numStr string

	for n > 0 {
		rem := n % 10
		numStr = string(rune('0'+rem)) + numStr
		n /= 10
	}

	if isnegative {
		numStr = "-" + numStr
	}

	return numStr
}
