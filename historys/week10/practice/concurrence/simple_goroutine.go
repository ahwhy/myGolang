package main

import (
	"fmt"
	"time"
)

func Add(a, b int) int {
	fmt.Println("Add")
	return a + b
}

var add = func(a, b int) int {
	fmt.Println("add")
	return a + b
}

func main2() {
	go Add(2, 4)
	go Add(3, 9)

	go func(a, b int) int {
		fmt.Println("add")
		return a + b
	}(2, 4)
	go func(a, b int) int {
		fmt.Println("add")
		return a + b
	}(3, 9)

	go add(2, 4)
	go add(3, 9)

	time.Sleep(10 * time.Millisecond)
}

//go run concurrence/simple_goroutine.go
