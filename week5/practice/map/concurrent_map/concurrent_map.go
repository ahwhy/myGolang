package main

import (
	"fmt"
	"time"
	cmap "github.com/orcaman/concurrent-map"
)

func main() {
	m := cmap.New()
	// 写map的 go
	go func() {
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			m.Set(key, i)
		}

	}()
	time.Sleep(5 * time.Second)
	// 读map的 go
	go func() {
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			v, exists := m.Get(key)
			if exists {
				fmt.Println(v.(int), exists)
			}
		}
	}()

	time.Sleep(100 * time.Second)

}
