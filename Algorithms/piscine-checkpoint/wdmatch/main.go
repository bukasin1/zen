package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		return
	}
	str1 := args[0]
	str2 := args[1]
	// fmt.Println("args:", str1, str2)

	if len(str1) > len(str2) {
		return
	}

	i := 0
	count := 0
	for _, ch := range str1 {
		// fmt.Println("str1 ch:", string(ch), i)
		for i < len(str2) {
			if rune(str2[i]) == ch {
				count++
				i++
				break
			}
			i++
		}
		if count == len(str1) {
			fmt.Println(str1)
			return
		}
		if i == len(str2) {
			return
		}
	}
	if count == len(str1) {
		fmt.Println(str1)
	}
}
