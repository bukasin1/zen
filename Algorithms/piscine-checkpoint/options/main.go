package main

import (
	"fmt"
	"os"
)

func printOptionsHelp() {
	fmt.Println("options: abcdefghijklmnopqrstuvwxyz")
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		printOptionsHelp()
		return
	}
	var intbits = make([]int, 32)
	// fmt.Println("initial:", intbits)

	for _, option := range args {
		if option[0] != '-' || len(option) == 1 {
			fmt.Println("Invalid Option")
			return
		}

		for i, ch := range option[1:] {
			if i == 0 && ch == 'h' {
				printOptionsHelp()
				return
			}
			if ch < 'a' || ch > 'z' {
				fmt.Println("Invalid Option")
				return
			}
			bitIndex := 31 - int(ch-'a')
			intbits[bitIndex] = 1
		}
	}

	for i, bit := range intbits {
		fmt.Print(string(rune('0' + bit)))
		if (i+1)%8 == 0 && i != len(intbits)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}
