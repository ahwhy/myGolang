package main_test

import (
	"container/list"
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

func TestListBasic(t *testing.T) {
	lst := list.New()
	lst.PushBack(1)
	lst.PushBack(2)
	lst.PushBack(3)
	lst.PushFront(4)
	lst.PushFront(5)
	lst.PushFront(6)
	TraversList(lst)
	ReverseList(lst)

	_ = lst.Remove(lst.Back())  //移除元素的同时返回元素的值，注意元素不能是nil
	_ = lst.Remove(lst.Front()) //Remove操作复杂度为O(1)
	fmt.Printf("length %d\n", lst.Len())
	TraversList(lst)
}

//go run data_structure/list.go
