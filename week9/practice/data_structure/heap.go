package main

import (
	"container/heap"
	"fmt"
)

type Item struct {
	Value    string
	priority int //优先级，数字越大，优先级越高
	index    int //item在堆里的编号
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority //golang默认提供的是小根堆，而优先队列是大根堆，所以这里要反着定义Less
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

//往slice里append,需要传slice指针
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

//让slice指向新的子切片，需要传slice指针
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1] //数组最后一个元素
	item.index = n - 1
	*pq = old[0 : n-1] //去掉最一个元素
	return item
}

func main5() {
	pq := make(PriorityQueue, 0, 10)
	pq.Push(&Item{"A", 3, 0}) //往数组里面添加元素
	pq.Push(&Item{"B", 2, 1})
	pq.Push(&Item{"C", 4, 2})
	heap.Init(&pq)                   //根据数组中的元素构建堆
	heap.Push(&pq, &Item{"D", 6, 3}) //通过heap添加元素
	for pq.Len() > 0 {
		fmt.Println(heap.Pop(&pq)) //通过heap删除堆顶元素
	}
}
