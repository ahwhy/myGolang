package main

import "fmt"

func append1(pq *PriorityQueue) {
	*pq = append(*pq, &Item{"A", 3, 0})
}

func append2(pq PriorityQueue) {
	pq = append(pq, &Item{"A", 3, 0}) //按值传递，并不会修改函数外面的pq
}

func slice1(pq *PriorityQueue) {
	n := len(*pq)
	old := *pq
	*pq = old[0 : n-1]
}

func slice2(pq PriorityQueue) {
	n := len(pq)
	old := pq
	pq = old[0 : n-1] //按值传递，并不会修改函数外面的pq
}

func main6() {
	pq := make(PriorityQueue, 0, 10)
	pq = append(pq, &Item{"D", 6, 3})
	append2(pq)
	for _, ele := range pq {
		fmt.Println(ele)
	}
	fmt.Println("=============")
	append1(&pq)
	for _, ele := range pq {
		fmt.Println(ele)
	}
	fmt.Println("=============")
	slice2(pq)
	for _, ele := range pq {
		fmt.Println(ele)
	}
	fmt.Println("=============")
	slice1(&pq)
	for _, ele := range pq {
		fmt.Println(ele)
	}
	fmt.Println("=============")
}

//go run data_structure/slice_pointer.go data_structure/heap.go
