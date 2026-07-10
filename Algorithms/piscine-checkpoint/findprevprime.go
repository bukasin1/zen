package main

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

func FindPrevPrime(nb int) int {
	for nb > 1 {
		if isPrime(nb) {
			return nb
		}
		nb--
	}

	return 0
}
