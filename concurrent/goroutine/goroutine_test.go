package goroutine_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

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


func TestOnce(t *testing.T){
	go goroutine.LoadResource()
	go goroutine.LoadResource()

	inst1 := goroutine.GetSingletonInstance()
	inst2 := goroutine.GetSingletonInstance()

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("inst1 address %v\n", []*goroutine.Singleton{inst1})
	fmt.Printf("inst2 address %v\n", []*goroutine.Singleton{inst2})
} 

func TestRuntime(t *testing.T) {
	fmt.Printf("逻辑处理器数目:%d\n", runtime.NumCPU())
	fmt.Printf("NumGoroutine:%d\n", runtime.NumGoroutine())
	fmt.Printf("NumCgoCall:%d\n", runtime.NumCgoCall())
	fmt.Printf("GOROOT:%s\n", runtime.GOROOT())
}
