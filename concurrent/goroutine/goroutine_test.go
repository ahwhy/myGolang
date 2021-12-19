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

func TestOnce(t *testing.T) {
	go goroutine.LoadResource()
	go goroutine.LoadResource()

	inst1 := goroutine.GetSingletonInstance()
	inst2 := goroutine.GetSingletonInstance()

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("inst1 address %v\n", []*goroutine.Singleton{inst1})
	fmt.Printf("inst2 address %v\n", []*goroutine.Singleton{inst2})
}

func TestLimit(t *testing.T) {
	go func() {
		ticker := time.NewTicker(1 * time.Second) //每隔1秒打印一次协程数量
		for {
			<-ticker.C
			fmt.Printf("当前协程数: %d\n", runtime.NumGoroutine())
		}
	}()

	work := func() {
		//do something
		time.Sleep(100 * time.Millisecond)
	}
	glimit := goroutine.NewGlimit(10) //限制协程数为10
	for i := 0; i < 10000; i++ {
		glimit.Run(work) //不停地通过Run创建子协程
	}
	time.Sleep(10 * time.Second)
}

func TestRuntime(t *testing.T) {
	fmt.Printf("逻辑处理器数目:%d\n", runtime.NumCPU())
	fmt.Printf("当前协程数:%d\n", runtime.NumGoroutine())
	fmt.Printf("NumCgoCall:%d\n", runtime.NumCgoCall())
	fmt.Printf("GOROOT:%s\n", runtime.GOROOT())
}

func TestApplication(t *testing.T) {
	goroutine.Application()
}

func TestDealPanic(t *testing.T) {
	goroutine.DealPanic()
}

func TestDealPanicInG(t *testing.T) {
	goroutine.DealPanicInG()
}

func TestDealPanicInGV2(t *testing.T) {
	goroutine.DealPanicInGV2()
}
