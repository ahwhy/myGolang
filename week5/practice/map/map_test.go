package main

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestMap_init(t *testing.T) {
	m := make(map[string]int)
	keys := make([]string, 0)
	// 填充数据到map
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key_%d", i)
		keys = append(keys, key)
		m[key] = i
	}
	fmt.Println(m)
	// range遍历 无序的
	for k, v := range m {
		fmt.Printf("[%s=%d]\n", k, v)
	}

	// 按key升序打印map
	fmt.Println("按key升序打印map")
	fmt.Println(keys)
	for _, k := range keys {
		v := m[k]
		fmt.Printf("[%s=%d]\n", k, v)
	}

	// 使用sort包排序后打印map
	sort.Strings(keys)
	fmt.Println(keys)
	for _, k := range keys {
		v := m[k]
		fmt.Printf("[%s=%d]\n", k, v)
	}
}

func TestMap_float(t *testing.T) {
	m := make(map[float64]int)
	m[2.4] = 2
	fmt.Printf("k: %v, v: %d\n", 2.4000000000000000000000001, m[2.4000000000000000000000001])
}

func TestMap_doublemap(t *testing.T) {
	doubleMap := make(map[string]map[string]string)
	v1 := make(map[string]string)
	v1["k1"] = "v1"
	doubleMap["k1"] = v1
	fmt.Println(doubleMap)
}

func TestMap_rw(t *testing.T) {
	c := make(map[int]int)
	// 写map的go
	go func() {
		for i := 0; i < 10000; i++ {
			c[i] = i
		}
	}()
	// 读map的go
	go func() {
		for i := 0; i < 10000; i++ {
			c[i] = i
			fmt.Println(c[i])
		}
	}()
	/*
		fatal error: concurrent map writes

		goroutine 5 [running]:
		runtime.throw(0xc2b26b, 0x15)
		        C:/Program Files/Go/src/runtime/panic.go:1117 +0x79 fp=0xc000047f60 sp=0
		xc000047f30 pc=0xbd2819

	*/
	fmt.Println(c)
	time.Sleep(20 * time.Second)
}
