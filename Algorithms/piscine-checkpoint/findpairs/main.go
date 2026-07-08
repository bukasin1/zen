package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Invalid input1.")
		return
	}

	arrStr := args[0]
	// err := json.Unmarshal([]byte(arrStr), &arr)
	if arrStr[0] != '[' || arrStr[len(arrStr)-1] != ']' {
		fmt.Println("Invalid input.")
		return
	}

	targetStr := args[1]
	targetNum, err1 := strconv.Atoi(targetStr)
	if err1 != nil {
		fmt.Println("Invalid target sum.")
		return
	}

	//form array slice
	var strArr []string

	// strArr = strings.Fields(arrStr[1 : len(arrStr)-1])

	var str string
	for i, ch := range arrStr {
		if i > 0 && i < len(arrStr)-1 {
			if ch == ' ' {
				strArr = append(strArr, str)
				str = ""
				continue
			} else if ch != ',' {
				str += string(ch)
			}
		}
	}
	strArr = append(strArr, str)

	var arr = make([]int, 0, len(strArr))
	for _, num := range strArr {
		if num[len(num)-1] == ',' {
			num = num[:len(num)-1]
		}
		n, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Invalid number:", num)
			return
		}
		arr = append(arr, n)
	}

	// fmt.Println(arr, targetNum, strArr)

	var out [][2]int
	for i, n := range arr {
		for j := i + 1; j < len(arr); j++ {
			p := arr[j]
			if n+p == targetNum && i != j {
				out = append(out, [2]int{i, j})
			}
		}
	}

	if len(out) == 0 {
		fmt.Println("No pairs found.")
	} else {
		fmt.Printf("Pairs with sum %d: %v\n", targetNum, out)
	}
}
