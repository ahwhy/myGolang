package goroutibe

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func printChars(prefix string, group *sync.WaitGroup) {
	defer group.Done()
	for ch := 'A'; ch <= 'Z'; ch++ {
		fmt.Printf("%s:%c\n", prefix, ch)
		runtime.Gosched()
	}
}

func TestSyncGroup(t *testing.T) {
	group := &sync.WaitGroup{}
	n := 10

	group.Add(n)
	for i := 0; i < n; i++ {
		go printChars(fmt.Sprintf("gochars%0d\n", i), group)
	}
	group.Wait()
	fmt.Println("over")
}

func TestCloser(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("-----------------------")
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
