package channel

import (
	"fmt"
	"sync"
)

// 协程A和B交替执行
func SyncAB() {
	send, recv := make(chan struct{}), make(chan struct{})
	wg := sync.WaitGroup{}

	wg.Add(2)
	go A(send, recv, wg)
	go B(recv, send, wg)

	recv <- struct{}{}
	wg.Wait()
}

func A(send, recv chan struct{}, wg sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	msg := []string{"1", "2", "3"}
	index := 0

	for range recv {
		if index > 2 {
			return
		}

		fmt.Println(msg[index])
		index++
		send <- struct{}{}
	}
}

func B(send, recv chan struct{}, wg sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	msg := []string{"X", "Y", "Z"}
	index := 0

	for range recv {
		if index > 2 {
			return
		}

		fmt.Println(msg[index])
		index++
		send <- struct{}{}
	}
}
