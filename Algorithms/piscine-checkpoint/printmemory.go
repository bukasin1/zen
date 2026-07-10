package main

import (
	"github.com/01-edu/z01"
)

func printString(s string) {
	for _, ch := range s {
		z01.PrintRune(ch)
	}
	z01.PrintRune('\n')
}

func PrintMemory(arr [10]byte) {
	var hexex = "0123456789abcdef"
	var str string
	for i, b := range arr {

		var mem []rune

		rem := b % 16

		// if rem > 9 {
		// 	rem = rem - 10 + 'a'
		// } else {
		// 	rem = rem + '0'
		// }
		rem = hexex[rem]
		mem = append(mem, rune(rem))

		hex := b / 16
		for hex > 0 {
			if hex >= 16 {
				rem = hex % 16
				if rem > 9 {
					rem = rem - 10 + 'a'
				} else {
					rem = rem + '0'
				}
				mem = append(mem, rune(rem))
			} else {
				var val byte
				if hex > 9 {
					val = hex - 10 + 'a'
				} else {
					val = hex + '0'
				}
				mem = append(mem, rune(val))
				// break
			}
			hex = hex / 16
		}

		// fmt.Println("mem:", mem)

		for i := len(mem) - 1; i >= 0; i-- {
			z01.PrintRune(mem[i])
		}
		// z01.PrintRune(rune(hex))
		// z01.PrintRune(rune(rem))
		if (i+1)%4 == 0 || i == len(arr)-1 {
			z01.PrintRune('\n')
		} else {
			z01.PrintRune(' ')
		}

		if b < 32 || b > 126 {
			str += "."
		} else {
			str += string(b)
		}
	}
	printString(str)
}
