package main

import (
	"fmt"
	"os"
	"time"
)

//倒计时
func countDown(countCh chan int, n int, finishCh chan struct{}) {
	if n <= 0 { //从n开始倒数
		return
	}
	ticker := time.NewTicker(1 * time.Second) //创建一个周期性的定时器，每隔1秒执行一次
	for {
		countCh <- n //把n放入管道
		<-ticker.C   //等1秒钟
		n--          //n减1
		if n <= 0 {  //n减到0时退出
			ticker.Stop()          //停止定时器
			finishCh <- struct{}{} //成功结束
			break                  //退出for循环
		}
	}
}

//中止
func abort(ch chan struct{}) {
	buffer := make([]byte, 1)
	os.Stdin.Read(buffer) //阻塞式IO，如果标准输入里没数据，该行一直阻塞。注意在键盘上敲完后要按下Enter才会把输入发给Stdin
	ch <- struct{}{}
}

func main15() {
	countCh := make(chan int)
	finishCh := make(chan struct{})
	go countDown(countCh, 10, finishCh) //开一个子协程，去往countCh和finishCh里放数据
	abortCh := make(chan struct{})
	go abort(abortCh) //开一个子协程，去往abortCh里放数据

LOOP:
	for { //循环监听
		select { //同时监听3个channel，谁先准备好就执行谁，然后进入下一次for循环
		case n := <-countCh:
			fmt.Println(n)
		case <-finishCh:
			fmt.Println("finish")
			break LOOP //退出for循环。在使用for select时，单独一个break不能退出for循环
		case <-abortCh:
			fmt.Println("abort")
			break LOOP //退出for循环
		}
	}
}
