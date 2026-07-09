package main

import (
	"fmt"
	"os"
)

func isVowel(ch rune) bool {
	vowels := "aeiou"
	for _, v := range vowels {
		// fmt.Println("vowel check:", string(v), string(v-' '), i)
		if ch == v || ch == (v-' ') {
			return true
		}
	}
	return false
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		return
	}

	var cons string

	for i, ch := range args[0] {
		if isVowel(ch) {
			pigword := args[0][i:] + cons + "ay"
			fmt.Println(pigword)
			return
		}
		cons += string(ch)
	}
	fmt.Println("No vowels")
}
