package main

import (
	"fmt"
	"sync"
	"time"
)

type concurrentMap struct {
	sync.RWMutex
	mp map[int]int
}

func (c *concurrentMap) Set(key int, value int) {
	// 先获取写锁
	c.Lock()
	// set值
	c.mp[key] = value
	// 解锁
	c.Unlock()
}
func (c *concurrentMap) Get(key int) int {
	// 先获取读锁
	c.RLock()
	// 获取值
	res := c.mp[key]
	// 解锁
	c.RUnlock()
	return res
}

func main() {
	// fatal error: concurrent map read and map write
	c := concurrentMap{
		mp: (map[int]int{}),
	}
	// 写map的goroutine
	go func() {
		for i := 0; i < 10000; i++ {

			c.Set(i, i)
		}
	}()
	// 读map的goroutine
	go func() {
		for i := 0; i < 10000; i++ {
			res := c.Get(i)
			fmt.Printf("[cmap.get][%d=%d]\n", i, res)
		}
	}()

	time.Sleep(20 * time.Second)
}
