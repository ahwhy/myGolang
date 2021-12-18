package selectio

import (
	"fmt"
	"time"
)

func TimeAfter() {
	fmt.Println(time.Now())

	a := time.After(6 * time.Second)
	fmt.Println(<-a)
}

func SelectTimeout() {
	ch1 := make(chan string)

	timeout := time.After(3 * time.Second)

	select {
	case val := <-ch1:
		fmt.Println("recv value from ch1:", val)
		return
	case val := <-timeout:
		fmt.Println(val)
		return
	}
}
