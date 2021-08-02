package main

import "fmt"

type ListNode struct {
	Value int
	Prev  *ListNode
	Next  *ListNode
}

type DoubleList struct {
	Head   *ListNode
	Tail   *ListNode
	Length int
}

func (list *DoubleList) Append(x int) {
	node := &ListNode{Value: x}
	tail := list.Tail
	if tail == nil {
		list.Head = node
		list.Tail = node
	} else {
		tail.Next = node
		node.Prev = tail
		list.Tail = node
	}
	list.Length += 1
}

func (list *DoubleList) Get(idx int) *ListNode {
	if list.Length <= idx {
		return nil
	}
	curr := list.Head
	for i := 0; i < idx; i++ {
		curr = curr.Next
	}
	return curr
}

func (list *DoubleList) InsertAfter(x int, prevNode *ListNode) {
	node := &ListNode{Value: x}
	if prevNode.Next == nil {
		prevNode.Next = node
		node.Prev = prevNode
	} else {
		nextNode := prevNode.Next
		nextNode.Prev = node
		node.Next = nextNode
		prevNode.Next = node
		node.Prev = prevNode
	}
}

func (list *DoubleList) Traverse() {
	curr := list.Head
	for curr != nil {
		fmt.Printf("%d ", curr.Value)
		curr = curr.Next
	}
	fmt.Println()
}

func main8() {
	list := new(DoubleList)
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Append(4)
	list.Append(5)
	list.Traverse()
	node := list.Get(3)
	list.InsertAfter(9, node)
	list.Traverse()
}
