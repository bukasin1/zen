package main

import (
	"fmt"
	"os"
	"strconv"
)

func printerror() {
	fmt.Println("Error")
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("len of args error")
		printerror()
		return
	}

	var operands []int
	var numstr string

	for i, op := range args[0] {
		if op == ' ' {
			if len(numstr) > 0 {
				num, err := strconv.Atoi(numstr)
				if err != nil {
					fmt.Println("invalid num error", numstr, string(op), i)
					printerror()
					return
				}
				operands = append(operands, num)
				numstr = ""
			}
			continue
		}
		switch op {
		case '+', '-', '*', '/':
			if len(operands) < 2 {
				fmt.Println("len of operands error", operands)
				printerror()
				return
			}
			var curans int
			switch op {
			case '+':
				curans = operands[0] + operands[1]
			case '-':
				curans = operands[0] - operands[1]
			case '*':
				curans = operands[0] * operands[1]
			case '/':
				curans = operands[0] / operands[1]
			}
			var remoperands []int
			if len(operands) > 2 {
				remoperands = operands[2:]
			}
			operands = append([]int{curans}, remoperands...)
			continue
		}
		numstr += string(op)
	}

	fmt.Println("Operands after operations:", operands)
	if len(operands) != 1 {
		printerror()
	} else {
		fmt.Println(operands[0])
	}
}
