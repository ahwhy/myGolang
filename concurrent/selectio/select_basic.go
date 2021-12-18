package selectio

import (
	"fmt"
	"time"
)

func Basic1() {
	ch1, ch2, ch3 := make(chan int, 1), make(chan int, 1), make(chan int, 1)
	ch1 <- 1
	ch2 <- 2
	ch3 <- 3

	select {
	case v := <-ch1:
		fmt.Println(v)
	case v := <-ch2:
		fmt.Println(v)
	case v := <-ch3:
		fmt.Println(v)
	default:
		fmt.Println("default")
	}
}

func Basic2() {
	ch1, ch2 := make(chan int, 5), make(chan int, 5)
	go pump1(ch1)
	go pump2(ch2)
	go suck(ch1, ch2)
	time.Sleep(1 * time.Second)
}

// 偶数channel
func pump1(ch chan int) {
	for i := 0; i <= 30; i++ {
		if i%2 == 0 {
			ch <- i
		}
	}
}

// 奇数channel
func pump2(ch chan int) {
	for i := 0; i <= 30; i++ {
		if i%2 == 1 {
			ch <- i
		}
	}
}

// 处理ch1和ch2中的数据
func suck(ch1 chan int, ch2 chan int) {
	for {
		select {
		case v := <-ch1:
			fmt.Printf("recv on ch1: %d\n", v)
		case v := <-ch2:
			fmt.Printf("recv on ch2: %d\n", v)
		}
	}
}
