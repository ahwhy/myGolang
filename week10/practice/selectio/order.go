package selectio

import "fmt"

func SelectOrder() {
	// 我们只在一个执行体(Goroutine)中完成
	// 所以需要channel是带缓冲的
	a := make(chan string, 100)
	b := make(chan string, 100)
	c := make(chan string, 100)

	// 依次往a, b, c 中放入数据
	for i := 0; i < 10; i++ {
		a <- "A"
		b <- "B"
		c <- "C"
	}

	
	for i := 0; i < 10; i++ {
		select {
		case v := <-a:
			fmt.Println(v)
		case v := <-b:
			fmt.Println(v)
		case v := <-c:
			fmt.Println(v)
		default:
			fmt.Println("Default")
		}
	}
}
