package stack

import (
	"sync"
)

// 定义需要存入的元素对象
// 这里Item是范型, 指代任意类型
type Item interface{}

type Stack_slice struct {
	items []Item
	lock  *sync.RWMutex
}

func NewStackslice() *Stack_slice {
	return &Stack_slice{
		items: []Item{},
		lock:  &sync.RWMutex{},
	}
}

func (stack *Stack_slice) Push(item interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.items = append(stack.items, item)
}

func (stack *Stack_slice) Pop() interface{} {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if stack.Empty() {
		return nil
	}
	item := stack.items[len(stack.items)-1]
	stack.items = stack.items[0 : len(stack.items)-1]
	return item
}

func (stack *Stack_slice) Peak() interface{} {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	return stack.items[len(stack.items)-1]
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

func (stack *Stack_slice) Clear() {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.items = []Item{}
}

// ForEach 遍历栈
func (stack *Stack_slice) ForEach(fn func(interface{})) {
	for i := stack.Len() - 1; i >= 0; i-- {
		fn(stack.items[i])
	}
}

// Sort 插入排序 把stack的元素从大到小进行排序 
func (s *Stack_slice) Sort() {
	// 准备一个辅助栈
	orderedStack := NewStackslice()

	for !s.Empty() {
		// 取出了顶层元素
		current := s.Pop()

		// 放入辅助栈 进行比较, 左边大于右边的比如(1 > 0), 交互该元素的位置
		for !orderedStack.Empty() && current.(int) > orderedStack.Peak().(int) {
			s.Push(orderedStack.Pop())
		}

		// 直接放入
		orderedStack.Push(current)
	}

	//  要把数据倒过来
	for !orderedStack.Empty() {
		s.Push(orderedStack.Pop())
	}
}

func NewNumberStack(numbers []int) *Stack_slice {
	items := make([]Item, 0, len(numbers))
	for i := range numbers {
		items = append(items, numbers[i])
	}
	return &Stack_slice{
		items: items,
	}
}
