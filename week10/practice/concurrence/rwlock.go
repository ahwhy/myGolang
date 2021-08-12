package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var n int32 = 0
var lock sync.RWMutex

func inc1() {
	n++ //n++不是原子操作，它分为3步：取出n，加1，结果赋给n
}

func inc2() {
	atomic.AddInt32(&n, 1) //封装成原子操作
}

func inc3() {
	lock.Lock()   //加写锁
	n++           //任一时刻，只有一个协程能进入临界区域
	lock.Unlock() //释放写锁
}

func main13() {
	const P = 1000 //开大量协程才能把脏写问题测出来
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc1()
		}()
	}
	wg.Wait()
	fmt.Printf("finally n=%d\n", n) //多运行几次，n经常不等于1000
	fmt.Println("===========================")
	n = 0 //重置n
	wg = sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc2()
		}()
	}
	wg.Wait()
	fmt.Printf("finally n=%d\n", atomic.LoadInt32(&n))
	fmt.Println("===========================")
	n = 0 //重置n
	wg = sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc3()
		}()
	}
	wg.Wait()
	lock.RLock() //加读锁。当写锁被其他协程持有时，加读锁操作将被阻塞；否则，如果其他协程持有读锁，加读锁操作不会被阻塞
	fmt.Printf("finally n=%d\n", n)
	lock.RUnlock() //释放读锁
	fmt.Println("===========================")
}

//go run concurrence/rwlock.go
//go run -race  concurrence/rwlock.go 	检测资源竞争
