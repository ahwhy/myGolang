package stack

import (
	"container/list"
	"fmt"
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

func (stack *Stack_list) Push(value interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.list.PushBack(value)
}

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

func (stack *Stack_list) Clear() {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.list.Init()
}

// ForEach 遍历栈
func (stack *Stack_list) ForEach(fn func(interface{})) {
	stack.lock.RLock()
	defer stack.lock.RUnlock()
	top := stack.list.Back()
	for top.Next() != nil {
		fmt.Printf("%v\n",top.Value)
		top = top.Next()
	}
}