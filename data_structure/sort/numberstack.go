package sort

import "fmt"

// 定义需要存入的元素对象
// 这里Item是范型, 指代任意类型
type Item interface{}

// 构建函数
func NewStack() *Stack {
	return &Stack{
		items: []Item{},
	}
}

type Stack struct {
	items []Item
}

func NewNumberStack(numbers []int) *Stack {
	items := make([]Item, 0, len(numbers))
	for i := range numbers {
		items = append(items, numbers[i])
	}
	return &Stack{
		items: items,
	}
}

// Push adds an Item to the top of the stack
func (s *Stack) Push(item Item) {
	s.items = append(s.items, item)
}

// Pop removes an Item from the top of the stack
func (s *Stack) Pop() Item {
	if s.IsEmpty() {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	return item
}

// Len 栈的大小
func (s *Stack) Len() int {
	return len(s.items)
}

// IsEmpty 判断是否为空
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Peek 获取栈顶元素的值 Peek
func (s *Stack) Peek() Item {
	if s.IsEmpty() {
		return nil
	}
	return s.items[len(s.items)-1]
}

// Clear 清空栈
func (s *Stack) Clear() {
	s.items = []Item{}
}

// Search 查询某个值 距离栈顶的距离
func (s *Stack) Search(item Item) (pos int, err error) {
	for i := range s.items {
		if item == s.items[i] {
			return i, nil
		}
	}
	return 0, fmt.Errorf("item %s not found", item)
}

// 遍历栈 ForEach
func (s *Stack) ForEach(fn func(Item)) {
	for i := range s.items {
		fn(i)
	}
}

// Sort 插入排序 把stack的元素从大到小进行排序 插入排序
func (s *Stack) Sort() {
	// 准备一个辅助的stack, 另一个容器
	orderdStack := NewStack()

	for !s.IsEmpty() {
		// 然后开始的排序流程
		current := s.Pop()

		// orderdStack顶端大于current，应该将orderdStack顶端移至s，直到orderdStack顶端小于current
		for !orderdStack.IsEmpty() && current.(int) > orderdStack.Peek().(int) {
			s.Push(orderdStack.Pop())
		}

		// 此时 当前current 一定是 <= orderdStack顶端
		orderdStack.Push(current)
	}

	// 倒过来
	for !orderdStack.IsEmpty() {
		s.Push(orderdStack.Pop())
	}
}
