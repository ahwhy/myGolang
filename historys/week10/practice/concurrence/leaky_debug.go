package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main18() {
	go func() {
		//在8080端口接收debug
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		//每隔1秒打印一次协程数量
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine())
		}
	}()

	for { //模拟在线服务，不停地处理请求，永不退出
		handle()
	}
}

//go run concurrence/leaky_debug.go concurrence/routine_leaky.go
//http://127.0.0.1:8080/debug/pprof/goroutine?debug=1
