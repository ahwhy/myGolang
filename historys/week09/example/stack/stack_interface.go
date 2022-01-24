package stack

type Stool interface {
	Clear()
	Len() int
	Empty() bool

	// Push 将元素压入栈
	Push(value interface{})

	// Pop 弹出栈最上方元素
	Pop() interface{}

	// Peak 返回栈最上方元素
	Peak() interface{}

	// ForEach 遍历栈
	ForEach(fn func(interface{}))

	// Sort 插入排序 把stack的元素从大到小进行排序
	Sort()
}
