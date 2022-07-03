package list_test

import (
	"container/list"
	"container/ring"
	"fmt"
	"testing"
)

func TraversList(lst *list.List) {
	head := lst.Front()
	for head.Next() != nil {
		fmt.Printf("%v ", head.Value)
		head = head.Next()
	}

	fmt.Println(head.Value)
}

func ReverseList(lst *list.List) {
	tail := lst.Back()
	for tail.Prev() != nil {
		fmt.Printf("%v ", tail.Value)
		tail = tail.Prev()
	}

	fmt.Println(tail.Value)
}

func TraverseRing(ring *ring.Ring) {
	ring.Do(func(i interface{}) { // 通过Do()来遍历ring，内部实际上调用了 Next() 而非 Prev()
		fmt.Printf("%v \n", i)
	})
}

func TestListPkgBasic(t *testing.T) {
	lst := list.New()
	lst.PushBack(1)
	lst.PushBack(2)
	lst.PushBack(3)
	lst.PushFront(4)
	lst.PushFront(5)
	lst.PushFront(6)
	TraversList(lst)
	ReverseList(lst)

	_ = lst.Remove(lst.Back())  // 移除元素的同时返回元素的值，注意元素不能是nil
	_ = lst.Remove(lst.Front()) // Remove操作复杂度为O(1)
	fmt.Printf("length %d\n", lst.Len())
	TraversList(lst)
}

func TestRingPkgBasic(t *testing.T) {
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
