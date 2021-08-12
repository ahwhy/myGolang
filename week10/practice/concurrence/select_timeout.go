package main

import (
	"context"
	"fmt"
	"time"
)

//模拟一个耗时较长的任务
func LongTimeWork() {
	time.Sleep(time.Duration(500) * time.Millisecond)
	return
}

//模拟一个接口处理函数
func Handle() {
	//借助于带超时的context来实现对函数的超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100) //改成1000试试
	defer cancel()                                                                 //纯粹出于良好习惯，函数退出前调用cancel()
	workDone := make(chan struct{}, 1)
	go func() { //把要控制超时的函数放到一个协程里
		LongTimeWork()
		workDone <- struct{}{}
	}()
	select { //下面的case只执行最早到来的那一个
	case <-workDone:
		fmt.Println("LongTimeWork return")
	case <-ctx.Done(): //ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		fmt.Println("LongTimeWork timeout")
	}
}

//模拟一个接口处理函数。超时控制的另一种实现
func Handle2() {
	//通过显式sleep再调用cancle()来实现对函数的超时控制
	ctx, cancel := context.WithCancel(context.Background())

	workDone := make(chan struct{}, 1)
	go func() { //把要控制超时的函数放到一个协程里
		LongTimeWork()
		workDone <- struct{}{}
	}()

	go func() {
		//100毫秒后调用cancel()，关闭ctx.Done()
		time.Sleep(100 * time.Millisecond) //改成1000试试
		cancel()
	}()

	select { //下面的case只执行最早到来的那一个
	case <-workDone:
		fmt.Println("LongTimeWork return")
	case <-ctx.Done(): //ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		fmt.Println("LongTimeWork timeout")
	}
}

func main16() {
	Handle()
	Handle2()
}

//go run concurrence/select_timeout.go
