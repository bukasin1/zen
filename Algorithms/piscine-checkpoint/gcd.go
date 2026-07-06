package main

func Gcd(a, b uint) uint {
	if a == 0 || b == 0 {
		return 0
	}

	var div uint = 1
	var lowest uint
	if a < b {
		lowest = a
	} else {
		lowest = b
	}

	for i := range lowest {
		// fmt.Println("range n:", i)
		if i > 0 && a%i == 0 && b%i == 0 {
			div = i
		}
	}

	return div
}
