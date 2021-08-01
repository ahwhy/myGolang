package stack

func NewNumberStack(numbers []int) *Stack {
	items := make([]Item, 0, len(numbers))
	for i := range numbers {
		items = append(items, numbers[i])
	}
	return &Stack{
		store: items,
	}
}

// 范型
type Item interface{}

// 构建函数
func NewStack() *Stack {
	return &Stack{}
}

type Stack struct {
	store []Item
}

func (s *Stack) Len() int {
	return len(s.store)
}

// 出栈, 弹出 去除切片尾部元素, 然后删除
func (s *Stack) Pop() Item {
	if s.IsEmpty() {
		return nil
	}

	//  [1 2] 2    index 1
	tail := s.Peek()

	// [1 2] [:1] [)
	s.store = s.store[:s.Len()-1]
	return tail
}

// 入栈，压栈
func (s *Stack) Push(item Item) {
	s.store = append(s.store, item)
}

// 空的栈
func (s Stack) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Stack) Peek() Item {
	if s.IsEmpty() {
		return nil
	}

	return s.store[s.Len()-1]
}

func (s *Stack) Clear() {
	s.store = []Item{}
}

// js 遍历
func (s *Stack) ForEach(fn func(Item)) {
	// 4 3 2 1 0
	for i := s.Len() - 1; i >= 0; i-- {
		fn(s.store[i])
	}
}

// 排序的方法
func (s Stack) Sort() {
	// 准备一个辅助栈
	orderedStack := NewStack()

	for !s.IsEmpty() {
		// 取出了顶层元素
		current := s.Pop()

		// 放入辅助栈 进行比较, 左边大于右边的比如(1 > 0), 交互该元素的位置
		for !orderedStack.IsEmpty() && current.(int) > orderedStack.Peek().(int) {
			s.Push(orderedStack.Pop())
		}

		// 直接放入
		orderedStack.Push(current)
	}

	//  要把数据倒过来
	for !orderedStack.IsEmpty() {
		s.Push(orderedStack.Pop())
	}
}
