package main

import (
	"fmt"
)

func sum(a int, b int) int {
	return a + b
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	fmt.Printf("sum 1 + 2: %d\n", sum(1, 2))
	fmt.Printf("fibonacci 20: %d\n", fibonacci(20))
}
