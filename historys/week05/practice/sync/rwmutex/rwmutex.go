package main

import (
	"log"
	"sync"
	"time"
)

var HcMutex sync.Mutex
var rwMutex sync.RWMutex

func runMutex(id int) {
	log.Printf("[任务id:%d][尝试获取锁]", id)
	HcMutex.Lock()
	log.Printf("[任务id:%d][获取到了锁][开始干活:睡眠10秒]", id)
	time.Sleep(10 * time.Second)
	HcMutex.Unlock()
	log.Printf("[任务id:%d][干完活了 释放锁]", id)
}
func runHcLock() {
	go runMutex(1)
	go runMutex(2)
	go runMutex(3)
}

func runReadLock(id int) {
	log.Printf("[读任务:%d][进入读方法等待获取读锁]", id)
	rwMutex.RLock()
	log.Printf("[读任务:%d][获取到了读锁][干活:睡眠10秒]", id)
	time.Sleep(10 * time.Second)
	rwMutex.RUnlock()
	log.Printf("[读任务:%d][完成读任务，释放读锁]", id)
}

func runWriteLock(id int) {
	log.Printf("[写任务:%d][进入写方法等待获取写锁]", id)
	rwMutex.Lock()
	log.Printf("[写任务:%d][获取到了写锁][干活:睡眠10秒]", id)
	time.Sleep(10 * time.Second)
	rwMutex.Unlock()
	log.Printf("[写任务:%d][完成写任务，释放写锁]", id)
}

func allWriteWorks() {
	for i := 1; i < 6; i++ {
		go runWriteLock(i)
	}
}

func allReadWorks() {
	for i := 1; i < 6; i++ {
		go runReadLock(i)
	}
}

func writeFirst() {
	go runWriteLock(1)
	time.Sleep(1 * time.Second)
	go runReadLock(1)
	go runReadLock(2)
	go runReadLock(3)
	go runReadLock(4)
	go runReadLock(5)
}

func readFirst() {
	go runReadLock(1)
	go runReadLock(2)
	go runReadLock(3)
	go runReadLock(4)
	go runReadLock(5)
	time.Sleep(1 * time.Second)
	go runWriteLock(1)
}

func main() {
	// 执行互斥锁 效果任务
	runHcLock()

	// 同时多个写锁任务，说明如果 并非使用 读写锁的写锁任务时，退化成了互斥锁
	allWriteWorks()

	// 同时并发 读锁任务， 说明读写锁的读锁 是可以同时施加多个读锁的
	allReadWorks()

	// 先启动写锁任务，后并非5个读锁任务，当有写锁时，读锁施加不了，写锁释放完，读锁可以施加多个
	writeFirst()

	// 并发5个读锁任务，启动一个写锁任务， 当有读锁时，阻塞写锁
	readFirst()

	time.Sleep(6000 * time.Second)
}
