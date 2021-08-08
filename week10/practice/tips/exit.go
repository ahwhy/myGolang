package tips

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(cannel chan struct{}) {
	for {
		select {
		default:
			fmt.Println("hello")
			time.Sleep(100 * time.Millisecond)
		case <-cannel:
			// 退出
			return
		}
	}
}

func CancelWithChannel() {
	cancel := make(chan struct{})
	go worker(cancel)

	time.Sleep(time.Second)
	cancel <- struct{}{}
}

func workerv2(wg *sync.WaitGroup, cancel chan bool) {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
			time.Sleep(100 * time.Millisecond)
		case <-cancel:
			// 清理工作需要进行
			return
		}
	}
}

func CancelWithDown() {
	cancel := make(chan bool)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go workerv2(&wg, cancel)
	}

	time.Sleep(time.Second)

	// 发送退出信号
	close(cancel)

	// 等待goroutine 安全退出
	wg.Wait()
}

func workerV3(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
			time.Sleep(100 * time.Millisecond)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
func CancelWithCtx() {
	start := time.Now()

	defer func() {
		fmt.Print(time.Since(start).Seconds())
	}()

	// 控制超时
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go workerV3(ctx, &wg)
	}

	// 等待安全退出
	wg.Wait()
}
