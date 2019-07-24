package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter count to get its fibonacci or q to quit\n")
	for {
		s, _ := reader.ReadString('\n')
		s = strings.Trim(s, "\n")
		if s == "" {
			continue
		}
		if s == "q" {
			return
		}
		count, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			fmt.Print("Entered count is not an integer\n")
			continue
		}
		go func(n int) {
			startedAt := time.Now()
			result := fibonacci(n)
			fmt.Printf("Fibonacci for %d is %d (took %d ms)\n", n, result, time.Now().Sub(startedAt)*time.Millisecond)
		}(int(count))
	}
}
