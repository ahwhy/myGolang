package goroutine_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ahwhy/myGolang/concurrent/goroutine"
)

var n int = 10

func TestSyncGroup(t *testing.T) {
	group := &sync.WaitGroup{}

	group.Add(n)
	for i := 0; i < n; i++ {
		go goroutine.PrintChars(fmt.Sprintf("gochars%0d\n", i), group)
	}

	group.Wait()
	fmt.Println("over")
}

// 闭包陷阱
func TestCloser(t *testing.T) {
	wg := &sync.WaitGroup{}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("-----------------------")

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
