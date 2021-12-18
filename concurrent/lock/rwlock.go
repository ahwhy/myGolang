package lock

import (
	"fmt"
	"sync"
	"time"
)

var (
	num    int64
	wg     sync.WaitGroup
	rwlock sync.RWMutex
)

func write() {
	// 加写锁
	rwlock.Lock()

	num = num + 1
	// 模拟真实写数据消耗的时间
	time.Sleep(10 * time.Millisecond)

	// 解写锁
	rwlock.Unlock()
	wg.Done()
}

func read() {
	// 加读锁
	rwlock.RLock()

	// 模拟真实读取数据消耗的时间
	time.Sleep(time.Millisecond)

	// 解读锁
	rwlock.RUnlock()

	// 退出协程前 记录 -1
	wg.Done()
}

func RWLock() {
	// 用于计算时间 消耗
	start := time.Now()

	// 开5个协程用作 写
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go write()
	}

	// 开500 个协程，用作读
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go read()
	}

	// 等待子协程退出
	wg.Wait()

	// 打印程序消耗的时间
	fmt.Println(time.Since(start))
}
