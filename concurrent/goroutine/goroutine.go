package goroutine

import (
	"fmt"
	"runtime"
	"sync"
)

func PrintChars(prefix string, group *sync.WaitGroup) {
	defer group.Wait()

	for ch := 'A'; ch <= 'Z'; ch++ {
		fmt.Printf("%s:%c\n", prefix, ch)
		runtime.Gosched()
	}
}
