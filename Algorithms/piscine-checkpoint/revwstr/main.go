package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		return
	}

	var words []string

	// words = strings.Fields(args[0])

	var word string
	for _, ch := range args[0] {
		if ch == ' ' {
			words = append(words, word)
			word = ""
			continue
		} else {
			word += string(ch)
		}
	}
	words = append(words, word)

	for i := len(words) - 1; i >= 0; i-- {
		fmt.Print(words[i])
		if i > 0 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}
