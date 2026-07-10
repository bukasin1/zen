package main

import (
	"fmt"
	// "github.com/01-edu/z01"
)

func printSolution(solution [8]int) {
	for i := 0; i < 8; i++ {
		fmt.Print(solution[i])
		// z01.PrintRune(rune(solution[i]) + '0')
	}
	fmt.Println()
	// z01.PrintRune('\n')
}

func isSafe(queens [8]int, col, row int) bool {
	for c := 0; c < col; c++ {
		if queens[c] == row || queens[c]-row == c-col || row-queens[c] == c-col {
			return false
		}
	}
	return true
}

func solve(queens *[8]int, col int) {
	// fmt.Println("solving queens:", queens, col)
	if col == 8 {
		printSolution(*queens)
		return
	}

	for row := 1; row <= 8; row++ {
		// fmt.Println("solving queens, row check:", queens, col, row)
		if isSafe(*queens, col, row) {
			queens[col] = row
			solve(queens, col+1)
		}
		// if row == 4 {
		// 	break
		// }
	}
}

func EightQueens() {
	var queens [8]int
	solve(&queens, 0)
	fmt.Println("Done solving queens!", queens)
}
