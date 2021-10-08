package main

import "fmt"

func main() {
	num := factorial(20)
	
	fmt.Println(num)
}

func factorial(n int) int {
	if n <= 0 {
		return -1
	} else if n == 1 {
		return 1
	} else {
		return n * factorial(n-1) // n * factorial(n-1) * factorial(n-2) * ... * 1  当n=1时，达到退出点
	}
}
