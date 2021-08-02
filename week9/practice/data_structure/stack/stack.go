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

func (stack *Stack) Push(value interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	node := &node{value, stack.top}
	stack.top = node
	stack.length++
}

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

func (stack *Stack) Peak() interface{} {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	if stack.length == 0 {
		return nil
	}
	return stack.top.value
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

func (stack *Stack) Clear() {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.top = nil
	stack.length = 0
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
