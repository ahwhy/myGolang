package list

import (
	"fmt"
)

type ListNode struct {
	Value interface{}
	Prev  *ListNode
	Next  *ListNode
}

type TList struct {
	Head   *ListNode
	Tail   *ListNode
	Length int
}

func (list *TList) Append(x interface{}) {
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

func (list *TList) Traverse() {
	curr := list.Head
	for curr != nil {
		fmt.Printf("%d ", curr.Value)
		curr = curr.Next
	}
	fmt.Println()
}

func (list *TList) Len() int {
	return list.Length
}

func (list *TList) Get(idx int) *ListNode {
	if list.Length <= idx {
		return nil
	}
	curr := list.Head
	for i := 0; i < idx; i++ {
		curr = curr.Next
	}
	return curr
}

func (list *TList) InsertAfter(x interface{}, prevNode *ListNode) {
	// prevNode -- node -- nextNode
	node := &ListNode{Value: x}
	if prevNode.Next == nil {
		prevNode.Next = node
		node.Prev = prevNode
		list.Tail = node
	} else {
		nextNode := prevNode.Next
		nextNode.Prev = node
		node.Next = nextNode
		prevNode.Next = node
		node.Prev = prevNode
	}
	list.Length += 1
}

func (list *TList) InsertBefore(x interface{}, nextNode *ListNode) {
	// prevNode -- node -- nextNode
	node := &ListNode{Value: x}
	if nextNode.Prev == nil {
		nextNode.Prev = node
		node.Next = nextNode
		list.Head = node
	} else {
		prevNode := nextNode.Prev
		nextNode.Prev = node
		node.Next = nextNode
		prevNode.Next = node
		node.Prev = prevNode
	}
	list.Length += 1
}

func (list *TList) Remove(idx int) *ListNode {
	if list.Length <= idx {
		return nil
	}
	curr := list.Head
	for i := 0; i < idx; i++ {
		curr = curr.Next
	}
	// prev --> curr --> next
	prev := curr.Prev
	next := curr.Next
	prev.Next, next.Prev = next, prev
	return curr
}

func (l *TList) RemoveEm(current *ListNode) error {
	// prev --> current --> next
	prev := current.Prev
	next := current.Next
	prev.Next, next.Prev = next, prev
	return nil
}

// ChangeToRing 将链表头尾相连成环
func (l *TList) ChangeToRing() {
	head, tail := l.Head, l.Tail
	// head  -->  tail
	head.Prev = tail
	// head  <--  tail
	tail.Next = head
}

func (list *TList) TraverseRing() {
	curr := list.Head
	for i := 0; i < 20; i++ {
		fmt.Printf("%d ", curr.Value)
		curr = curr.Next
	}
	fmt.Println()
}
