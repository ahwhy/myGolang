package heap_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/data_structure/heap"
)

func TestHeapbig(t *testing.T) {
	m := []int{0, 9, 3, 6, 2, 1, 7} //第0个下标不放目标元素
	h := heap.NewIntHeap(m)
	fmt.Println(h.Items())

	h.Push(50)
	fmt.Println(h.Items())

	h.Pop()
	fmt.Println(h.Items())

	h.Pop()
	fmt.Println(h.Items())
}

func TestBuildHeap(t *testing.T) {
	arr := []int{62, 20, 30, 15, 10, 49, 78, 45, 12, 11, 45}
	heap.ReverseAdjust(arr)
	fmt.Println(arr)
}

func TestHeapcontainer(t *testing.T) {
	heap.Example_intHeap()
}

func TestPriorityQueue(t *testing.T) {
	heap.TestPriorityQueue()
}

func TestValuePoint(t *testing.T) {
	heap.TestValuePoint()
}
