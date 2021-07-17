package main

import (
	"log"
	"sync"
	"time"
)

var HcMutex sync.Mutex

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

func main() {
	// 执行互斥锁 效果任务
	runHcLock()

	time.Sleep(600 * time.Second)
}
