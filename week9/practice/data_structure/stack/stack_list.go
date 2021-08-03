package stack

import (
	"container/list"
	"sync"
)

// 借助list包，实现stack
type Stack_list struct {
	list *list.List
	lock *sync.RWMutex
}

func NewStacklist() *Stack_list {
	list := list.New()
	lock := &sync.RWMutex{}
	return &Stack_list{list, lock}
}

func (stack *Stack_list) Clear() {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.list.Init()
}

func (stack *Stack_list) Len() int {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	return stack.list.Len()
}

func (stack *Stack_list) Empty() bool {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	return stack.list.Len() == 0
}

// Push 将元素压入栈
func (stack *Stack_list) Push(value interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.list.PushBack(value)
}

// Pop 弹出栈最上方元素
func (stack *Stack_list) Pop() interface{} {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	ele := stack.list.Back()
	if ele == nil {
		return nil
	} else {
		return stack.list.Remove(ele) // 移除ele的同时返回ele的值，注意ele不能是nil
	}
}

// Peak 返回栈最上方元素
func (stack *Stack_list) Peak() interface{} {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	ele := stack.list.Back()
	if ele == nil {
		return nil
	} else {
		return ele.Value
	}
}

// ForEach 遍历栈
func (stack *Stack_list) ForEach(fn func(interface{})) {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	top := stack.list.Back()
	for top.Prev() != nil {
		fn(top.Value)
		top = top.Prev()
	}
	fn(top.Value)
}

// Sort 插入排序 把stack的元素从大到小进行排序
func (stack *Stack_list) Sort() {
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
