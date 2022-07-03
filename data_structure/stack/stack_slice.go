package stack

import (
	"sync"
)

// 定义需要存入的元素对象
// 这里Item是范型, 指代任意类型
type Item interface{}

func NewStackslice() *Stack_slice {
	return &Stack_slice{
		items: []Item{},
		lock:  &sync.RWMutex{},
	}
}

type Stack_slice struct {
	items []Item
	lock  *sync.RWMutex
}

func (stack *Stack_slice) Clear() {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	stack.items = []Item{}
}

func (stack *Stack_slice) Len() int {
	stack.lock.RLock()
	defer stack.lock.RUnlock()

	return len(stack.items)
}

func (stack *Stack_slice) Empty() bool {
	stack.lock.RLock()
	defer stack.lock.RUnlock()

	return len(stack.items) == 0
}

// Push 将元素压入栈
func (stack *Stack_slice) Push(item interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	stack.items = append(stack.items, item)
}

// Pop 弹出栈最上方元素
func (stack *Stack_slice) Pop() interface{} {
	// stack.lock.Lock()
	// defer stack.lock.Unlock()
	if stack.Empty() {
		return nil
	}

	item := stack.items[len(stack.items)-1]
	stack.lock.Lock()
	stack.items = stack.items[0 : len(stack.items)-1]
	stack.lock.Unlock()

	return item
}

// Peak 返回栈最上方元素
func (stack *Stack_slice) Peak() interface{} {
	stack.lock.RLock()
	defer stack.lock.RUnlock()

	if stack.Empty() {
		return nil
	}

	return stack.items[len(stack.items)-1]
}

// ForEach 遍历栈
func (stack *Stack_slice) ForEach(fn func(interface{})) {
	for i := stack.Len() - 1; i >= 0; i-- {
		fn(stack.items[i])
	}
}

// Sort 插入排序 把stack的元素从大到小进行排序
func (stack *Stack_slice) Sort() {
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
