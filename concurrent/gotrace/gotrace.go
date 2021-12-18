package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

/*
Go语言中提供trace工具，用于分析程序的运行过程
 - 执行程序后，会生成trace.out文件
 - 再运行 go tool trace trace.out
*/

var wg sync.WaitGroup

func runTask(id int) {
	defer wg.Done()

	fmt.Printf("task %d strat ... \n", id)
	time.Sleep(1 * time.Second)
	fmt.Printf("task %d complete\n", id)
}

func asyncRun() {
	for i := 0; i < 10; i++ {
		go runTask(i + 1)
		wg.Add(1)
	}
}

func main() {
	// 创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	asyncRun()
	wg.Wait()
}
