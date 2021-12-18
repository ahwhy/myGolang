package channel

import "fmt"

func sender(ch chan string) {
	ch <- "hello"
	ch <- "this"
	ch <- "is"
	ch <- "alice"
	// 通话结束
	ch <- "EOF"

	close(ch)
}

// recver 循环读取chan里面的数据，直到channel关闭
func recver(ch chan string, down chan struct{}) {
	defer func() {
		down <- struct{}{}
	}()

	for v := range ch {
		if v == "EOF" {
			return
		}
		fmt.Println(v)
	}
}

func BasicChan() {
	ch := make(chan string)
	down := make(chan struct{})

	go sender(ch)       // sender goroutine
	go recver(ch, down) // recver goroutine

	<-down
}
