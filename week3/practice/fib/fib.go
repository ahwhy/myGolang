package main

import (
	"fmt"
)

func main() {
	for i := 2; i < 10; i++ {
		fmt.Printf("%d ", fib(i))
	}
}

func fib(n int) int {
	if n == 1 || n == 2 { // f(2) = f(1) = 1
		return 1
	}
	
	return fib(n-1) + fib(n-2) // f(n)=f(n-1)+f(n-2)
}
