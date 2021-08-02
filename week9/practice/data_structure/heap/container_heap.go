package heap

import (
	"container/heap"
	"fmt"
)

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int { return len(h) }

// 由于是小顶堆, 所以前面的元素(堆顶)要比后面的小
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }

// 交换2个元素的位置
func (h IntHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push 往堆里面添加元素, 后面的调整由heap包帮我们完成(shift-up)
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

// Pop 弹出堆顶的元素, 还记得堆是如何删除元素的吗?
// 1. 取出堆顶的元素给你
// 2. 那现在堆顶的空位应该被如何补充
// 3. 把堆底的元素给他, 然后就是shift-down
// 我们这里Pop 就是给他最后一个元素, 后面的操作由heap包的pop帮我们完成(shift-donw)
func (h *IntHeap) Pop() interface{} {
	// 保存old
	old := *h
	n := len(old)

	// 取出堆底部的元素x
	x := old[n-1]

	// 移除堆底的元素
	*h = old[0 : n-1]
	return x
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func Example_intHeap() {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minmum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
	// Output:
	// minimum: 1
	// 1 2 3 5
}
