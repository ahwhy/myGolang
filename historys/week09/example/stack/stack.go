package stack

import (
	"sync"
)

// 自定义链表，实现stack
type (
	node struct {
		value interface{}
		pev   *node
	}
	Stack struct {
		top    *node
		length int
		lock   *sync.RWMutex
	}
)

func NewStack() *Stack {
	return &Stack{nil, 0, &sync.RWMutex{}}
}

func (stack *Stack) Clear() {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.top = nil
	stack.length = 0
}

func (stack *Stack) Len() int {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	return stack.length
}

func (stack *Stack) Empty() bool {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	return stack.length == 0
}

// Push 将元素压入栈
func (stack *Stack) Push(value interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	node := &node{value, stack.top}
	stack.top = node
	stack.length++
}

// Pop 弹出栈最上方元素
func (stack *Stack) Pop() interface{} {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if stack.length == 0 {
		return nil
	}
	node := stack.top
	stack.top = node.pev
	stack.length--
	return node.value
}

// Peak 返回栈最上方元素
func (stack *Stack) Peak() interface{} {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	if stack.length == 0 {
		return nil
	}
	return stack.top.value
}

// ForEach 遍历栈
func (stack *Stack) ForEach(fn func(interface{})) {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	node := stack.top
	for node != nil {
		fn(node.value)
		node = node.pev
	}
}

// Sort 插入排序 把stack的元素从大到小进行排序
func (stack *Stack) Sort() {
	// 准备一个辅助的stack, 另一个容器
	orderedStack := NewStackslice()

	for !stack.Empty() {
		// 然后开始的排序流程 取出栈顶元素
		current := stack.Pop()

		// orderdStack顶端大于current，应该将orderdStack顶端移至stack，直到orderdStack顶端小于current
		for !orderedStack.Empty() && current.(int) > orderedStack.Peak().(int) {
			stack.Push(orderedStack.Pop())
		}

		// 直接放入
		orderedStack.Push(current)
	}

	//  要把数据倒过来
	for !orderedStack.Empty() {
		stack.Push(orderedStack.Pop())
	}
}
