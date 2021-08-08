package tips

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func DealPanic() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	arr := []int{0}
	_ = arr[2]
}

func DealPanicInG() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	wg.Add(1)
	go work()

	wg.Wait()
}

func work() {
	arr := []int{0}
	_ = arr[2]
	wg.Done()
}

func DealPanicInGV2() {
	// 处理主Goroutine的异常
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	wg.Add(1)
	go workV2()

	wg.Wait()
}

func workV2() {
	// 处理协程的异常
	defer func() {
		wg.Done()

		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	arr := []int{0}
	_ = arr[2]
}
