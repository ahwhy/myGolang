package main

import (
	"fmt"
	"sync"
)

/**
数组、slice、struct允许并发修改（可能会脏写），并发修改map会发生panic，因为map的value是不可寻址的。
如果需要并发修改map请使用sync.Map
*/

type Student struct {
	Name string
	Age  int32
}

var arr = [10]int{}
var m = sync.Map{}

func main14() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() { //写偶数位
		defer wg.Done()
		for i := 0; i < len(arr); i += 2 {
			arr[i] = 0
		}
	}()
	go func() { //写奇数位
		defer wg.Done()
		for i := 1; i < len(arr); i += 2 {
			arr[i] = 1
		}
	}()
	wg.Wait()
	fmt.Println(arr) //输出[0 1 0 1 0 1 0 1 0 1]
	fmt.Println("=======================")
	wg.Add(2)
	var stu Student
	go func() {
		defer wg.Done()
		stu.Name = "Fred"
	}()
	go func() {
		defer wg.Done()
		stu.Age = 20
	}()
	wg.Wait()
	fmt.Printf("%s %d\n", stu.Name, stu.Age)
	fmt.Println("=======================")
	wg.Add(2)
	go func() {
		defer wg.Done()
		m.Store("k1", "v1")
	}()
	go func() {
		defer wg.Done()
		m.Store("k1", "v2")
	}()
	wg.Wait()
	fmt.Println(m.Load("k1"))
}

//go run concurrence/collection_safety.go
