package selectio

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var (
	wg     sync.WaitGroup
	ctx    context.Context
	cancle context.CancelFunc
)

func Graceful_exit() {
	// 下面3个协程关联到了同一个context，通过cancle()可以通知彼此
	go listenSignal()
	go listenHttp(8080)
	go listenHttp(8081)

	wg.Wait() // 等待3个子协程优雅退出后，main协程再退出
}

func listenSignal() {
	defer wg.Done()

	c := make(chan os.Signal)
	//监听指定信号 SIGINT和SIGTERM。按下control+c向进程发送SIGINT信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done(): // 调用cancle() 管道ctx.Done()会被关闭，从ctx.Done()中读数据会立即返回0值
			return
		case sig := <-c:
			fmt.Printf("got signal %d\n", sig)
			cancle() // 取消 通知所有用到ctx的协程
			return
		}
	}
}

func listenHttp(port int) {
	defer wg.Done()

	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: nil} // 在端口port上开启http服务

	go func() {
		for {
			select {
			case <-ctx.Done():
				server.Close()
				return
			}
		}
	}()
	
	server.ListenAndServe() // 该行代码会一直阻塞，直到server.Close()
	fmt.Printf("stop listen on port %d\n", port)
}

func init() {
	wg = sync.WaitGroup{}
	wg.Add(3)                                              // 3个子协程，1个用于接收终止信号，其他2个是业务需要的后台协程
	ctx, cancle = context.WithCancel(context.Background()) // 父context
}
