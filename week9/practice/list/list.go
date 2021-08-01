package list

import (
	"fmt"
)

func NewIntNode(v int) *Node {
	return &Node{Value: v}
}

// 定义节点
type Node struct {
	// 存储你需要存储的数据
	Value interface{}
	// 下一跳的指向
	Next *Node
	// 上一跳
	Prev *Node
}

func NewIntList(headValue int) *List {
	head := &Node{Value: headValue}
	return &List{
		head: head,
	}
}

type List struct {
	head *Node
}

func (l *List) AddNode(n *Node) {
	// 我需要找到尾节点
	next := l.head
	for next.Next != nil {
		next = next.Next
	}

	// 修改为节点
	next.Next = n

	// 补充Previos指针
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

func (l *List) InsertAfter(after, current *Node) error {
	// 假设我们已经插入，他数据结构应该是啥样的
	// after --> current --> after_next

	// 保存下之前的after next
	afterNext := after.Next

	// 插入，修改指向
	after.Next = current
	current.Next = afterNext

	// 补充Previos指针
	// after <-- current <-- after_next
	current.Prev = after
	afterNext.Prev = current
	return nil
}

func (l *List) InsertBefore(current, n *Node) error {
	// 假设我们已经插入，他数据结构应该是啥样的
	//   -->   previous  -->     current -->
	//   previous --> n  --> current
	//

	// 保存下之前的before next
	previous := current.Prev

	// 插入，修改指向
	previous.Next = n
	n.Next = current

	// 补充Previos指针
	// before <-- current <-- before_next
	// current.Prev = before
	// beforeNext.Prev = current
	return nil
}

func (l *List) Remove(current *Node) error {
	// before --> current --> before_next
	prev := current.Prev
	prev.Next = current.Next
	return nil
}

// 变身成环
func (l *List) ChangeToRing() {
	// 我需要找到尾节点
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
