package main

import (
	"fmt"
	"time"
)

var cloch = make(chan int, 1)
var cloch2 = make(chan int, 1)

func traverseChannel() {
	for ele := range cloch {
		fmt.Printf("receive %d\n", ele)
	}
	fmt.Println()
}

func traverseChannel2() {
	for {
		if ele, ok := <-cloch2; ok { //ok==true代表管道还没有close
			fmt.Printf("receive %d\n", ele)
		} else { //管道关闭后，读操作会立即返回“0值”
			fmt.Printf("channel have been closed, receive %d\n", ele)
			break
		}
	}
}

func main11() {
	cloch <- 1
	close(cloch)
	traverseChannel() //如果不close就直接通过range遍历管道，会发生fatal error: all goroutines are asleep - deadlock!
	fmt.Println("==================")
	go traverseChannel2()
	cloch2 <- 1
	close(cloch2)
	time.Sleep(10 * time.Millisecond)
}

//go run concurrence/close_channel.go
