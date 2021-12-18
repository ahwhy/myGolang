package cspmodel

import (
	"fmt"
	"time"
)

func ProducerConsumerMode() {
	ch := make(chan int, 64) // 成果队列

	go Producer(3, ch) // 生成 3 的倍数的序列
	go Producer(5, ch) // 生成 5 的倍数的序列
	go Producer(7, ch) // 生成 5 的倍数的序列
	go Producer(9, ch) // 生成 5 的倍数的序列

	go Consumer(ch) // 消费 生成的队列
	go Consumer(ch) // 消费 生成的队列

	// 运行一定时间后退出
	time.Sleep(5 * time.Second)
}

// 生产者: 生成 factor 整数倍的序列， 3， 0, 3, 6, 9,  5, 0, 5, 10, 15
func Producer(factor int, out chan<- int) {
	maxCount := 0

	for i := 0; ; i++ {
		out <- i * factor

		// 最多生成10个
		maxCount++
		if maxCount > 10 {
			break
		}
	}
}

// 消费者
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}
