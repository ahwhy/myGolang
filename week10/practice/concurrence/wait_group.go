package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main3() { //main协程
	const N = 10
	wg := sync.WaitGroup{}
	wg.Add(N) //加N
	for i := 0; i < N; i++ {
		go func(a, b int) { //开N个子协程
			defer wg.Done() //减1
			time.Sleep(10 * time.Millisecond)
			_ = a + b
		}(i, i+1)
	}
	fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine())
	wg.Wait() //等待减为0
	fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine())
}

//go run concurrence/wait_group.go
