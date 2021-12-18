package main

import (
	"fmt"
	"time"
)

var syncChann = make(chan int)

func takeFromSyncChann() {
	if v, ok := <-syncChann; ok { //ok==true说明管道还没有关闭close
		fmt.Printf("take %d from synchronous channel\n", v)
	}
}

func putToSyncChann() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("putToSyncChann发生panic:%s\n", err)
		} else {
			fmt.Println("putToSyncChann success")
		}
	}()
	syncChann <- 10 //取操作没有准备好时，写管道会发生fatal error(注意不是panic，通过recover不能捕获）
}

func testsSyncChann1() {
	go takeFromSyncChann() //消费者协程会阻塞2秒钟，等待put
	time.Sleep(2 * time.Second)
	putToSyncChann()
}
func testsSyncChann2() {
	putToSyncChann() //取操作没有准备好时，写管道会发生fatal error
	go takeFromSyncChann()
}

func testsSyncChann3() {
	go func() {
		for {
			if v, ok := <-syncChann; ok { //ok==true说明管道还没有关闭close
				fmt.Printf("receive %d\n", v)
			} else {
				break
			}
		}
	}()
	for i := 0; i < 10; i++ {
		syncChann <- i
		fmt.Printf("send %d\n", i)
	}
	close(syncChann)
}

func main7() {
	// testsSyncChann1()
	// time.Sleep(3 * time.Second)
	// fmt.Println("==========================")
	// testsSyncChann3()
	// fmt.Println("==========================")

	testsSyncChann2()
	time.Sleep(3 * time.Second)
	fmt.Println("==========================")
}

//go run concurrence/sync_channel.go
