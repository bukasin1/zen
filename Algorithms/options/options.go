package main

import (
	"fmt"
	"os"
)

func main() {
	optionArgs := os.Args[1:]

	if len(optionArgs) == 0 || optionArgs[0] == "-h" {
		fmt.Println("options: abcdefghijklmnopqrstuvwxyz")
		return
	}

	output := make([]bool, 32)
	var nilByte byte

	fmt.Println("output:", output, nilByte)

	for _, option := range optionArgs {
		if len(option) <= 1 || option[0] != '-' {
			fmt.Println("Invalid Option")
			return
		}

		if option[1] == 'h' {
			fmt.Println("options: abcdefghijklmnopqrstuvwxyz")
			return
		}

		for _, c := range option[1:] {
			if c >= 'a' && c <= 'z' {
				alphaIndex := int(c - 'a' + 1)
				bit := 32 - alphaIndex
				output[bit] = true
			}
		}
	}

	// fmt.Println("output after:", output)
	result := ""

	for i, bitBool := range output {
		bit := "0"
		if bitBool {
			bit = "1"
		}
		result += string(bit)
		if (i+1)%8 == 0 {
			result += " "
		}
	}

	fmt.Println("output string:", result)
}
