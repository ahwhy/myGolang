package main

import "fmt"

func fib() func() int { // 斐波那契数列
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {
	fun := fib()
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
}
