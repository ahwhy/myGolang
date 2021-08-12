package main

import (
	"fmt"
	"time"
)

func main4() {
	arr := []int{1, 2, 3, 4}
	for _, v := range arr {
		go func() {
			fmt.Printf("%d\t", v) //用的是协程外面的全局变量v。输出4 4 4 4
		}()
	}
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println()
	for _, v := range arr {
		go func(value int) {
			fmt.Printf("%d\t", value) //输出1 4 2 3
		}(v) //把v的副本传到协程内部
	}
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println()
}

//go run concurrence/closure.go
