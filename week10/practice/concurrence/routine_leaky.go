package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

//模拟一个耗时较长的任务
func work() {
	time.Sleep(time.Duration(500) * time.Millisecond)
	return
}

//模拟一个接口处理函数
func handle() {
	//借助于带超时的context来实现对函数的超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100) //改成1000试试
	defer cancel()
	// begin := time.Now()             //纯粹出于良好习惯，函数退出前调用cancel()
	workDone := make(chan struct{}) //创建一个无缓冲管道
	go func() {                     //启动一个子协程
		work()
		workDone <- struct{}{} //work()结束后到，走到这行代码会一直阻塞，子协程无法结束，导致协程泄漏
	}()
	select { //下面的case只执行最早到来的那一个
	case <-workDone: //永远执行不到
		fmt.Println("LongTimeWork return")
	case <-ctx.Done(): //ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		// fmt.Printf("LongTimeWork timeout %d ms\n", time.Since(begin).Milliseconds())
	}
}

func main17() {
	for i := 0; i < 10; i++ {
		handle()
	}
	time.Sleep(2 * time.Second)                      //等所有work()结束
	fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine()) //11，10个阻塞的子协程 加 main协程
}

//go run concurrence/select_timeout.go
