package stack

type Stool interface {
	Clear()
	Len() int
	Empty() bool
	Push(value interface{})
	Pop() interface{}
	Peak() interface{}
	ForEach(fn func(interface{}))
}
