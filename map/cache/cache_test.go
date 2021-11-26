package cache_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ahwhy/myGolang/map/cache"
)

func TestCache(t *testing.T) {
	c := cache.NewCache()
	// 异步执行，删除过期项目的任务
	go c.Gc(10)

	// 写入数据
	for i := 0; i <= 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		item := cache.NewItem(i)

		// 设置缓存
		log.Printf("[设置项目缓存[key:%s][v:%v]", key, item)
		c.Set(key, item)
	}

	time.Sleep(15 * time.Second)

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		item := cache.NewItem(i + 2)

		// 更新缓存
		log.Printf("[更新项目缓存][key:%s][v:%v]", key, item)
		c.Set(key, item)
	}

	time.Sleep(10 * time.Second)
}
