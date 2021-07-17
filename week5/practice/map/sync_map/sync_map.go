package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

func main() {
	m := sync.Map{}
	// 增加
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		m.Store(key, i)
	}

	// 读取
	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("key_%d", i)
		value, exists := m.Load(key)
		if exists {
			v := value.(int)
			log.Printf("[%s=%d]", key, v)
		} else {
			log.Printf("[key_not_found_in_sync.map:%s]", key)
		}
	}

	// 删除
	m.Delete("key_9")

	// Range遍历
	m.Range(func(k, v interface{}) bool {
		key := k.(string)
		value := v.(int)
		log.Printf("[找到了][%s=%d]", key, value)
		return true
	})
	log.Printf("return false的遍历")
	m.Range(func(k, v interface{}) bool {
		key := k.(string)
		if strings.HasSuffix(key, "3") {
			log.Printf("不想要3")
			return false
		} else {
			log.Printf("[找到了][%s=%d]", key, v.(int))
			return true
		}
	})
	// LoadOrStore
	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("key_%d", i)
		v, loaded := m.LoadOrStore(key, i)
		if loaded {
			//说明之前有
			v := v.(int)
			log.Printf("[LoadOrStore][之前有][%s=%d]", key, v)
		} else {
			// 之前没有，新添加的
			log.Printf("[LoadOrStore][之前没有，新添加的][%s=%d]", key, v.(int))
		}
	}

	value, loaded := m.LoadAndDelete("key_11")
	log.Printf("[LoadAndDelete][key_11][v:%v][loaded:%v]", value, loaded)
}
