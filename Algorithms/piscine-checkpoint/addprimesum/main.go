package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/01-edu/z01"
)

func isPrime(n int) bool {
	if n == 2 || n == 3 || n == 5 || n == 7 {
		return true
	}

	if n < 2 || n%2 == 0 || n%3 == 0 || n%5 == 0 || n%7 == 0 {
		return false
	}

	for i := 11; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		z01.PrintRune('0')
		z01.PrintRune('\n')
		return
	}
	numStr := args[0]
	n, e := strconv.Atoi(numStr)
	if e != nil || n < 0 {
		z01.PrintRune('0')
		z01.PrintRune('\n')
		return
	}

	var sum = 0
	for n > 0 {
		if isPrime(n) {
			sum += n
		}
		n--
	}

	fmt.Println(sum)
}
