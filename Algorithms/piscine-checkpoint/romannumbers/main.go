package main

import (
	"fmt"
	"os"
	"strconv"
)

func printError() {
	fmt.Println("ERROR: cannot convert to roman digit")
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		printError()
		return
	}

	decNum, err := strconv.Atoi(args[0])

	if err != nil || decNum >= 4000 {
		printError()
		return
	}

	// romdigits := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romannumeralsMap := map[int]string{
		1:    "I",
		4:    "IV",
		5:    "V",
		9:    "IX",
		10:   "X",
		40:   "XL",
		50:   "L",
		90:   "XC",
		100:  "C",
		400:  "CD",
		500:  "D",
		900:  "CM",
		1000: "M",
	}

	var romnum, romcalc string

	// for decNum > 0 {
	// 	for _, d := range romdigits {
	// 		if decNum >= d {
	// 			var romdigit = romannumeralsMap[d]
	// 			calcdigit := romdigit
	// 			if len(calcdigit) == 2 {
	// 				calcdigit = "(" + string(calcdigit[1]) + "-" + string(calcdigit[0]) + ")"
	// 			}
	// 			if len(romcalc) == 0 {
	// 				romcalc += calcdigit
	// 			} else {
	// 				romcalc += "+" + calcdigit
	// 			}
	// 			romnum += romdigit
	// 			decNum -= d
	// 			break
	// 		}
	// 	}
	// }

	for decNum > 0 {
		var times, dec int
		if decNum >= 1000 {
			dec = 1000
		} else if decNum >= 900 {
			dec = 900
		} else if decNum >= 500 {
			dec = 500
		} else if decNum >= 400 {
			dec = 400
		} else if decNum >= 100 {
			dec = 100
		} else if decNum >= 90 {
			dec = 90
		} else if decNum >= 50 {
			dec = 50
		} else if decNum >= 40 {
			dec = 40
		} else if decNum >= 10 {
			dec = 10
		} else if decNum >= 9 {
			dec = 9
		} else if decNum >= 5 {
			dec = 5
		} else if decNum >= 4 {
			dec = 4
		} else if decNum >= 1 {
			dec = 1
		}
		times = decNum / dec
		decNum = decNum % dec

		var romdigit = romannumeralsMap[dec]
		for range times {
			calcdigit := romdigit
			if len(calcdigit) == 2 {
				calcdigit = "(" + string(calcdigit[1]) + "-" + string(calcdigit[0]) + ")"
			}
			if len(romcalc) == 0 {
				romcalc += calcdigit
			} else {
				romcalc += "+" + calcdigit
			}
			romnum += romdigit
		}
	}

	fmt.Println(romcalc)
	fmt.Println(romnum)
}
