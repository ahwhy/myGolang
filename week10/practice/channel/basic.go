package chnanel

import (
	"fmt"
)

func sender(ch chan string) {
	ch <- "hello"
	ch <- "this"
	ch <- "is"
	ch <- "alice"
	// 发送通话结束
	ch <- "EOF"
	close(ch)
}

// recver 循环读取chan里面的数据，直到channel关闭
func recver(ch chan string, down chan struct{}) {
	defer func() {
		down <- struct{}{}
	}()

	for v := range ch {
		// 处理通话结束
		if v == "EOF" {
			return
		}
		fmt.Println(v)
	}
}

func Basic() {
	ch := make(chan string)

	down := make(chan struct{}) // bool string struct{}{}
	go sender(ch)               // sender goroutine
	go recver(ch, down)         // recver goroutine

	<-down
}
