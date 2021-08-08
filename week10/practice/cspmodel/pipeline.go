package cspmodel

import (
	"fmt"
	"math/rand"
	"sync"
)

var wg sync.WaitGroup

func PipelineMode() {
	wg.Add(3)

	// 创建两个channel
	ch1 := make(chan int)
	ch2 := make(chan int)

	// 3个goroutine并行
	go getRandNum(ch1)
	go addRandNum(ch1, ch2)
	go printRes(ch2)

	wg.Wait()
}

func getRandNum(out chan<- int) {
	// defer the wg.Done()
	defer wg.Done()

	var random int
	// 总共生成10个随机数
	for i := 0; i < 10; i++ {
		// 生成[0,30)之间的随机整数并放进channel out
		random = rand.Intn(30)
		out <- random
	}
	close(out)
}

func addRandNum(in <-chan int, out chan<- int) {
	defer wg.Done()
	for v := range in {
		// 输出从第一个channel中读取到的数据
		// 并将值+1后放进第二个channel中
		fmt.Println("before +1:", v)
		out <- (v + 1)
	}
	close(out)
}

func printRes(in <-chan int) {
	defer wg.Done()
	for v := range in {
		fmt.Println("after +1:", v)
	}
}
