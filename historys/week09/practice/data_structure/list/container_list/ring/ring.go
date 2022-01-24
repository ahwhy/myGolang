package main

import (
	"container/ring"
	"fmt"
)

func TraverseRing(ring *ring.Ring) {
	ring.Do(func(i interface{}) { //通过Do()来遍历ring，内部实际上调用了Next()而非Prev()
		fmt.Printf("%v ", i)
	})
	fmt.Println()
}

func main() {
	ring := ring.New(5) //必须指定长度，各元素被初始化为nil
	ring2 := ring.Prev()
	for i := 0; i < 3; i++ {
		ring.Value = i
		ring = ring.Next()
	}
	for i := 0; i < 3; i++ {
		ring2.Value = i
		ring2 = ring2.Prev()
	}
	TraverseRing(ring)
	TraverseRing(ring2) //ring和ring2当前所在的指针位置不同，所以遍历出来的顺序也不同
}

//go run data_structure/ring.go
