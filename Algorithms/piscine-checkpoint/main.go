package main

import "fmt"

func main() {
	fmt.Println("-----CheckNumber----")
	fmt.Println(CheckNumber("Hello"))
	fmt.Println(CheckNumber("Hello1"))

	fmt.Println("\n-----CountAlpha------")
	fmt.Println(CountAlpha("Hello world"))
	fmt.Println(CountAlpha("H e l l o"))
	fmt.Println(CountAlpha("H1e2l3l4o"))

	fmt.Println("\n---RetainFirstHalf---")
	fmt.Println(RetainFirstHalf("This is the 1st halfThis is the 2nd half"))
	fmt.Println(RetainFirstHalf("A"))
	fmt.Println(RetainFirstHalf(""))
	fmt.Println(RetainFirstHalf("Hello World"))

	fmt.Println("\n---CamelToSnakeCase--")
	fmt.Println(CamelToSnakeCase("HelloWorld"))
	fmt.Println(CamelToSnakeCase("helloWorld"))
	fmt.Println(CamelToSnakeCase("camelCase"))
	fmt.Println(CamelToSnakeCase("CAMELtoSnackCASE"))
	fmt.Println(CamelToSnakeCase("camelToSnakeCase"))
	fmt.Println(CamelToSnakeCase("CcCc"))

	fmt.Println("\n-----FirstWord----")
	fmt.Print(FirstWord("hello there"))
	fmt.Print(FirstWord("   "))
	fmt.Print(FirstWord(" hello   .........  bye"))

	fmt.Println("\n-----Gcd----")
	fmt.Println(Gcd(42, 10))
	fmt.Println(Gcd(42, 12))
	fmt.Println(Gcd(14, 77))
	fmt.Println(Gcd(17, 3))

	fmt.Println("\n-----HashCode----")
	fmt.Println(HashCode("A"))
	fmt.Println(HashCode("AB"))
	fmt.Println(HashCode("BAC"))
	fmt.Println(HashCode("Hello World"))

	fmt.Println("\n-----FindPrevPrime----")
	fmt.Println(FindPrevPrime(6))
	fmt.Println(FindPrevPrime(4))

	fmt.Println("\n-----PrintMemory----")
	PrintMemory([10]byte{'h', 'e', 'l', 'l', 'o', 255, 21, '*'})

	// The character '👋' (wave emoji) takes 4 bytes in UTF-8
	// msg := "Hi👋!"

	// fmt.Println("--- Looping by Bytes (index access) ---")
	// for i := 0; i < len(msg); i++ {
	// 	fmt.Printf("Byte %d: %d (as char: %c)\n", i, msg[i], msg[i])
	// }

	// fmt.Println("\n--- Looping by Runes (for range) ---")
	// for index, runeValue := range msg {
	// 	fmt.Printf("Rune starts at byte %d: %d (as char: %c)\n", index, runeValue, runeValue)
	// }

	fmt.Println("\n-----ConcatSlice----")
	fmt.Println(ConcatSlice([]int{1, 2, 3}, []int{4, 5, 6}))
	fmt.Println(ConcatSlice([]int{}, []int{4, 5, 6, 7, 8, 9}))
	fmt.Println(ConcatSlice([]int{1, 2, 3}, []int{}))

	fmt.Println("\n-----ConcatAlternate----")
	fmt.Println(ConcatAlternate([]int{1, 2, 3}, []int{4, 5, 6}))
	fmt.Println(ConcatAlternate([]int{2, 4, 6, 8, 10}, []int{1, 3, 5, 7, 9, 11}))
	fmt.Println(ConcatAlternate([]int{1, 2, 3}, []int{4, 5, 6, 7, 8, 9}))
	fmt.Println(ConcatAlternate([]int{1, 2, 3}, []int{}))

	fmt.Println("\n-----RevConcatAlternate----")
	fmt.Println(RevConcatAlternate([]int{1, 2, 3}, []int{4, 5, 6}))
	fmt.Println(RevConcatAlternate([]int{1, 2, 3}, []int{4, 5, 6, 7, 8, 9}))
	fmt.Println(RevConcatAlternate([]int{1, 2, 3, 9, 8}, []int{4, 5}))
	fmt.Println(RevConcatAlternate([]int{1, 2, 3}, []int{}))

	fmt.Println("\n-----Itoa----")
	fmt.Println(Itoa(12345))
	fmt.Println(Itoa(0))
	fmt.Println(Itoa(-1234))
	fmt.Println(Itoa(987654321))
	fmt.Println(Itoa(-9223372036854775808))

	fmt.Println("\n-----ItoaBase----")
	fmt.Println(ItoaBase(10, 2))
	fmt.Println(ItoaBase(255, 16))
	fmt.Println(ItoaBase(-42, 4))
	fmt.Println(ItoaBase(123, 10))
	fmt.Println(ItoaBase(0, 8))
	fmt.Println(ItoaBase(255, 2))
	fmt.Println(ItoaBase(-255, 16))
	fmt.Println(ItoaBase(15, 16))
	fmt.Println(ItoaBase(10, 4))
	fmt.Println(ItoaBase(255, 10))
}
