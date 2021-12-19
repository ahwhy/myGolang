package selectio

import (
	"fmt"
	"os"
	"time"
)

func SelectChannel() {
	countCh := make(chan int)
	finishCh := make(chan struct{})
	abortCh := make(chan struct{})

	go countDown(countCh, 10, finishCh)
	go abort(abortCh)

LOOP:
	for {
		select { //同时监听3个channel，谁先准备好就执行谁，然后进入下一次for循环
		case n := <-countCh:
			fmt.Println(n)
		case <-finishCh:
			fmt.Println("finish")
			break LOOP //退出for循环 在使用for select时，单独一个break不能退出for循环
		case <-abortCh:
			fmt.Println("abort")
			break LOOP
		}
	}
}

// 倒计时
func countDown(countCh chan int, n int, finishCh chan struct{}) {
	//  从n开始倒数
	if n < 0 {
		return
	}

	// 创建一个周期性的定时器，每隔1秒执行一次
	ticker := time.NewTicker(1 * time.Second)
	for {
		countCh <- n
		<-ticker.C // 等待1s
		n--
		if n <= 0 {
			ticker.Stop()
			finishCh <- struct{}{}
			break
		}
	}
}

// 中止
func abort(ch chan struct{}) {
	buffer := make([]byte, 10)

	// 阻塞式IO，如果标准输入里没数据，该行一直阻塞
	// 需要在键盘上敲完后要按下Enter才会把输入发给Stdin
	os.Stderr.Read(buffer)
	ch <- struct{}{}
}
