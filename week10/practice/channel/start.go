package chnanel

import (
	"fmt"
)

func A(send chan struct{}, recv chan struct{}) {
	defer wg.Done()

	msg := []string{"1", "2", "3"}

	count := 0
	for range recv {
		if count > 2 {
			return
		}

		fmt.Println(msg[count])
		count++
		send <- struct{}{}
	}
}

func B(send chan struct{}, recv chan struct{}) {
	defer wg.Done()

	msg := []string{"A", "B", "C"}

	count := 0
	for range recv {
		if count > 2 {
			return
		}

		fmt.Println(msg[count])
		count++
		send <- struct{}{}
	}

}

func SyncAB() {
	a, b := make(chan struct{}), make(chan struct{})

	wg.Add(2)
	go A(a, b)
	go B(b, a)

	// 启动信号
	b <- struct{}{}

	wg.Wait()
}
