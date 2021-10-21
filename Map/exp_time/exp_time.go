package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type item struct {
	value int   // 值
	ts    int64 // 时间戳，item被创建出来的时间,或者被更新的时间
}

type Cache struct {
	sync.RWMutex
	mp map[string]*item
}

func (c *Cache) Get(key string) *item {
	c.RLock()
	defer c.RUnlock()
	return c.mp[key]
}

func (c *Cache) Set(key string, value *item) {
	c.Lock()
	defer c.Unlock()
	c.mp[key] = value
}

func (c *Cache) Gc(timeDelta int64) {
	// GC 先加读锁 -> 检查确实有需要回收的数据 -> 合并写锁回收。
	for {
		toDelKeys := make([]string, 0)
		now := time.Now().Unix()
		c.RLock()

		// 变量缓存中的项目，对比时间戳，超过 timeDelta的删除
		for k, v := range c.mp {
			if now-v.ts > timeDelta {
				log.Printf("[这个项目过期了][key %s]", k)
				toDelKeys = append(toDelKeys, k)
			}
		}
		c.RUnlock()

		c.Lock()
		for _, k := range toDelKeys {
			delete(c.mp, k)
		}
		c.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func main() {
	c := Cache{
		mp: make(map[string]*item),
	}
	// 让删除过期项目的任务，异步执行，
	go c.Gc(30)

	// 写入数据 从mysql读取
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		ts := time.Now().Unix()
		im := &item{
			value: i,
			ts:    ts,
		}
		//设置缓存
		log.Printf("[设置缓存][项目][key:%s][v:%v]", key, im)
		c.Set(key, im)
	}
	time.Sleep(31 * time.Second)
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		ts := time.Now().Unix()
		im := &item{
			value: i + 1,
			ts:    ts,
		}
		log.Printf("[更新缓存][项目][key:%s][v:%v]", key, im)
		c.Set(key, im)
	}
	select {} // 阻塞main
}
