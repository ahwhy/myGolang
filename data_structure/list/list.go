package list

import (
	"fmt"
)

func NewIntList(headValue interface{}) *List {
	// 链表的头
	head := &Node{Value: headValue}

	return &List{
		head: head,
	}
}

type List struct {
	head *Node
}

func NewIntNode(v interface{}) *Node {
	return &Node{Value: v}
}

// 重新定义节点
type Node struct {
	// 需要存储的数据
	Value interface{}
	// 下一跳
	Next *Node
	// 上一跳
	Prev *Node
}

func (l *List) AddNode(n *Node) {
	// 需要找到尾节点
	next := l.head
	for next.Next != nil {
		next = next.Next
	}

	// 修改为节点
	next.Next = n
	n.Prev = next
}

func (l *List) Traverse(fn func(n *Node)) {
	loopCount := 1

	n := l.head
	for n.Next != nil {
		// 最多循环5伦
		if loopCount > 5 {
			return
		}

		fn(n)
		n = n.Next

		if n == l.head {
			loopCount++
		}
	}

	fn(n)
	fmt.Println()
}

func (l *List) Len() int {
	len := 0
	n := l.head
	if n.Prev != nil {
		return -1
	}
	for n.Next != nil {
		n = n.Next
		len++
	}

	return len + 1
}

func (l *List) Get(idx int) interface{} {
	index := 0
	n := l.head
	for n.Next != nil {
		n = n.Next
		index++
		if index == idx {
			return n.Value
		}
	}

	return nil
}

func (l *List) InsertAfter(after, current *Node) error {
	// after --> current --> afterNext
	// 保存after的下一跳
	afterNext := after.Next

	// 插入current，修改指向
	after.Next = current
	current.Next = afterNext

	// after <-- current <-- after_next
	current.Prev = after
	afterNext.Prev = current

	return nil
}

func (l *List) InsertBefore(before, current *Node) error {
	// beforePrev <-- current <-- before
	// 保存before的上一跳
	beforePrev := before.Prev

	// 插入current，修改指向
	before.Prev = current
	current.Prev = beforePrev

	// beforePrev --> current --> before
	current.Next = before
	beforePrev.Next = current

	return nil
}

func (l *List) Remove(current *Node) error {
	// prev --> current --> next
	prev := current.Prev
	next := current.Next
	prev.Next, next.Prev = next, prev

	return nil
}

// ChangeToRing 将链表头尾相连成环
func (l *List) ChangeToRing() {
	// 需要找到尾节点
	next := l.head
	for next.Next != nil {
		next = next.Next
	}
	
	head, tail := l.head, next
	// head  -->  tail
	head.Prev = tail
	// head  <--  tail
	tail.Next = head
}
