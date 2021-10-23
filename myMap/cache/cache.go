package cache

import (
	"log"
	"sync"
	"time"
)

type item struct {
	value interface{} // 值
	ts    int64       // 时间戳，item被创建出来的时间，或者被更新的时间
}

func NewItem(i interface{}) *item {
	return &item{
		value: i,
		ts:    time.Now().Unix(),
	}
}

// 带过期缓存的map
type Cache struct {
	sync.RWMutex
	mp map[string]*item
}

func NewCache() *Cache {
	return &Cache{
		mp: make(map[string]*item),
	}
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

// GC 先加读锁 -> 检查确实有需要回收的数据 -> 合并写锁回收
func (c *Cache) Gc(timeDelta int64) {
	for {
		toDelKeys := make([]string, 0)
		now := time.Now().Unix()
		c.RLock()

		// 遍历缓存中的项目，对比时间戳，超过 timeDelta 就删除该项目
		for k, v := range c.mp {
			if now-v.ts > timeDelta {
				log.Printf("[项目已过期][key %s]", k)
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
