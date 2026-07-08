package main

func ConcatAlternate(slice1, slice2 []int) []int {
	var first []int
	var second []int

	if len(slice1) >= len(slice2) {
		first = slice1
		second = slice2
	} else {
		first = slice2
		second = slice1
	}

	// var out = make([]int, len(slice1)+len(slice2))
	// var lastIndex int

	var out = make([]int, 0, len(slice1)+len(slice2))

	for i, n := range first {
		// if i < len(second) {
		// 	out[i*2] = n
		// 	out[i*2+1] = second[i]
		// 	lastIndex = i*2 + 1 + 1
		// } else {
		// 	out[lastIndex] = n
		// 	lastIndex++
		// }

		out = append(out, n)
		if i < len(second) {
			out = append(out, second[i])
		}
	}

	return out
}
