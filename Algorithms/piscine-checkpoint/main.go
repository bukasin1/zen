package main

import "fmt"

func main() {
	fmt.Println(CheckNumber("Hello"))
	fmt.Println(CheckNumber("Hello1"))

	fmt.Println(CountAlpha("Hello world"))
	fmt.Println(CountAlpha("H e l l o"))
	fmt.Println(CountAlpha("H1e2l3l4o"))

	fmt.Println(RetainFirstHalf("This is the 1st halfThis is the 2nd half"))
	fmt.Println(RetainFirstHalf("A"))
	fmt.Println(RetainFirstHalf(""))
	fmt.Println(RetainFirstHalf("Hello World"))

	fmt.Println(CamelToSnakeCase("HelloWorld"))
	fmt.Println(CamelToSnakeCase("helloWorld"))
	fmt.Println(CamelToSnakeCase("camelCase"))
	fmt.Println(CamelToSnakeCase("CAMELtoSnackCASE"))
	fmt.Println(CamelToSnakeCase("camelToSnakeCase"))
	fmt.Println(CamelToSnakeCase("CcCc"))

	fmt.Print(FirstWord("hello there"))
	fmt.Print(FirstWord("   "))
	fmt.Print(FirstWord(" hello   .........  bye"))

	fmt.Println(Gcd(42, 10))
	fmt.Println(Gcd(42, 12))
	fmt.Println(Gcd(14, 77))
	fmt.Println(Gcd(17, 3))

	fmt.Println(HashCode("A"))
	fmt.Println(HashCode("AB"))
	fmt.Println(HashCode("BAC"))
	fmt.Println(HashCode("Hello World"))

	fmt.Println(FindPrevPrime(6))
	fmt.Println(FindPrevPrime(4))

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

}
