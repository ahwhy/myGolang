package lock

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	sumCount int64
	lock     sync.Mutex
)

// 普通版加函数
func add() {
	sumCount = sumCount + 1
	wg.Done()
}

// 互斥锁版加函数
func mutexAdd() {
	lock.Lock()
	sumCount = sumCount + 1
	lock.Unlock()
	wg.Done()
}

// 原子操作版加函数
func atomicAdd() {
	atomic.AddInt64(&sumCount, 1)
	wg.Done()
}

func ProblemV1() {
	var sum int

	// 使用WaitGroup等待1000个goroutine完成
	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			// 对变量count执行1000次加1
			for j := 0; j < 10; j++ {
				lock.Lock()
				sum++
				lock.Unlock()
			}
		}()
	}
	// 等待10个goroutine完成
	wg.Wait()

	// 结果出错，而且不稳定。
	fmt.Println(sum)
}

func ProblemV2() {
	// 目的是 记录程序消耗时间
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		// go mutexAdd() // 互斥锁的 add函数 是并发安全的，因为拿不到互斥锁会阻塞，所以加锁性能开销大
		go atomicAdd() // 原子操作的 add函数 是并发安全，性能优于加锁的
	}
	// 等待10个goroutine完成
	wg.Wait()
	fmt.Println(sumCount)
}
