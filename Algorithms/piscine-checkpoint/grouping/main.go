package main

import (
	"fmt"
	"os"
)

func parseregex(regex string) (parsed []string, ok bool) {
	if regex[0] != '(' && regex[len(regex)-1] != ')' {
		return parsed, false
	}

	var re string
	for _, ch := range regex[1 : len(regex)-1] {
		if ch == '|' {
			if len(re) > 0 {
				parsed = append(parsed, re)
			}
			re = ""
			continue
		}
		re += string(ch)
	}
	if len(re) > 0 {
		parsed = append(parsed, re)
	}

	if len(parsed) > 0 {
		return parsed, true
	}

	return parsed, ok
}

func printword(word string, count, currentnum int) {
	for i := range count {
		fmt.Printf("%d: %v\n", i+1+currentnum, word)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 || args[1] == "" {
		return
	}

	regexes, ok := parseregex(args[0])
	if !ok {
		return
	}
	fmt.Println("regexes:", regexes)

	var currentword string
	var count, currentnum int
	var seenregmap = make(map[string]bool)
	for i, ch := range args[1] {
		if ch == ' ' {
			if count > 0 && len(currentword) > 0 {
				printword(currentword, count, currentnum)
			}
			currentnum += count
			currentword = ""
			count = 0
			seenregmap = make(map[string]bool)
			continue
		}

		if ch != ',' {
			currentword += string(ch)
		}
		// if count < len(regexes) {
		for _, re := range regexes {
			relength := len(re)
			if i+relength <= len(args[1]) && !seenregmap[re] && re == args[1][i:i+relength] {
				seenregmap[re] = true
				count++
			}
		}
		// }
	}
	if count > 0 && len(currentword) > 0 {
		printword(currentword, count, currentnum)
	}
}
