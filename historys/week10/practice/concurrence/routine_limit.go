package main

import (
	"fmt"
	"runtime"
	"time"
)

type Glimit struct {
	limit int
	ch    chan struct{}
}

func NewGlimit(limit int) *Glimit {
	return &Glimit{
		limit: limit,
		ch:    make(chan struct{}, limit), //缓冲长度为limit，运行的协程不会超过这个值
	}
}

func (g *Glimit) Run(f func()) {
	g.ch <- struct{}{} //创建子协程前往管道里send一个数据
	go func() {
		f()
		<-g.ch //子协程退出时从管理里取出一个数据
	}()
}

func main19() {
	go func() {
		//每隔1秒打印一次协程数量
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine())
		}
	}()

	work := func() {
		//do something
		time.Sleep(100 * time.Millisecond)
	}
	glimit := NewGlimit(10) //限制协程数为10
	for i := 0; i < 10000; i++ {
		glimit.Run(work) //不停地通过Run创建子协程
	}
	time.Sleep(10 * time.Second)
}

//go run concurrence/routine_limit.go
