package channel

import "fmt"

func senderV2(ch chan string, down chan struct{}) {
	ch <- "hello"
	ch <- "this"
	ch <- "is"
	ch <- "alice"
	// 发送通话结束
	ch <- "EOF"

	// 同步模式等待recver 处理完成
	<-down
	
	// 处理完成后关闭channel
	close(ch)
}

// recver 循环读取chan里面的数据，直到channel关闭
func recverV2(ch chan string, down chan struct{}) {
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

func BufferedChan() {
	ch := make(chan string, 5)

	down := make(chan struct{})
	go senderV2(ch, down) // sender goroutine
	go recverV2(ch, down) // recver goroutine

	<-down
}
