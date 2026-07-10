package main

func RevConcatAlternate(slice1, slice2 []int) []int {
	var large, small []int
	var out = make([]int, 0, len(slice1)+len(slice2))

	if len(slice1) >= len(slice2) {
		large = slice1
		small = slice2
	} else {
		large = slice2
		small = slice1
	}

	for i := len(large) - 1; i >= 0; i-- {
		if i < len(small) {
			out = append(out, slice1[i])
			out = append(out, slice2[i])
		} else {
			out = append(out, large[i])
		}
	}

	return out
}
