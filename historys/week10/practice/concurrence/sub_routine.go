package main

import (
	"fmt"
	"runtime"
	"time"
)

func grandson() {
	fmt.Println("grandson begin")
	fmt.Printf("routine num %d\n", runtime.NumGoroutine())
	time.Sleep(8 * time.Second)
	fmt.Printf("routine num %d\n", runtime.NumGoroutine())
	fmt.Println("grandson finish")
}

func child() {
	fmt.Println("child begin")
	go grandson()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("child finish") //子协程退出后，孙协挰还在运行
}

func main23() {
	go child()
	time.Sleep(10 * time.Second)
}
