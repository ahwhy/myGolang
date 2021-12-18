package main

import (
	"fmt"
)

func moo(x int) int {
	fmt.Printf("x=%d\n", x)
	return x
}

func soo(a, b int) int {
	defer func() {
		//recover必须在defer中才能生效
		if err := recover(); err != nil { // err != nil说明发生了panic
			fmt.Printf("soo函数中发生了panic:%s\n", err)
		}
	}()
	c := a*3 + 9
	//defer是先进后出，即逆序执行
	defer fmt.Println("first defer")
	d := c + 5
	defer fmt.Println("second defer")
	e := d / b //如果发生panic，则后面的defer不会执行
	defer fmt.Println("third defer")
	return moo(e) //defer是在函数临退出前执行，不是在代码的return语句之前执行，因为return语句不是原子操作
}

func main6() {
	soo(5, 0)
}

//go run concurrence/recover.go
