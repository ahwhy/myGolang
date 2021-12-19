package selectio

import (
	"context"
	"fmt"
	"time"
)

func TimeAfter() {
	fmt.Println(time.Now())

	a := time.After(6 * time.Second)
	fmt.Println(<-a)
}

func SelectTimeout() {
	ch1 := make(chan string)

	timeout := time.After(3 * time.Second)

	select {
	case val := <-ch1:
		fmt.Println("recv value from ch1:", val)
		return
	case val := <-timeout:
		fmt.Println(val)
		return
	}
}

// 模拟一个耗时较长的任务
func longTimeWork() {
	time.Sleep(time.Duration(50) * time.Millisecond)
}

// 模拟一个接口处理函数
func handle() {
	//借助于带超时的context来实现对函数的超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*100)
	defer cancel()

	workDone := make(chan struct{}, 1)
	go func() { // 把要控制超时的函数放到一个协程里
		longTimeWork()
		workDone <- struct{}{}
	}()

	select { // 下面的case只执行最早到来的那一个
	case <-workDone:
		fmt.Println("LongTimeWork return")
	case <-ctx.Done(): // ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		fmt.Println("LongTimeWork timeout")
	}
}

// 模拟一个接口处理函数  超时控制的另一种实现
func handle2() {
	//借助于带超时的context来实现对函数的超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*1000)
	go func() {
		//1000毫秒后调用cancel()，关闭ctx.Done()
		time.Sleep(1000 * time.Millisecond) //改成1000试试
		cancel()
	}()

	workDone := make(chan struct{}, 1)
	go func() { // 把要控制超时的函数放到一个协程里
		longTimeWork()
		workDone <- struct{}{}
	}()

	select { // 下面的case只执行最早到来的那一个
	case <-workDone:
		fmt.Println("LongTimeWork return")
	case <-ctx.Done(): // ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		fmt.Println("LongTimeWork timeout")
	}
}

func SelectTimeoutV2() {
	handle()
	handle2()
}
